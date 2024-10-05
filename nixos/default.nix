flake: {
  pkgs,
  system,
  lib,
  config,
  ...
}: let
  website = flake.packages."${system}".default;
  cfg = config.services.website;
  serverOptions = with lib;
    types.submodule {
      options = {
        port = mkOption {
          type = types.number;
        };
      };
    };
  databaseOptions = with lib;
    types.submodule {
      options = {
        password = mkOption {
          type = types.str;
        };
      };
    };
  rootDatabaseOptions = with lib;
    types.submodule {
      options = {
        username = mkOption {
          type = types.str;
        };
        password = mkOption {
          type = types.str;
        };
      };
    };
  goaccessOptions = with lib;
    types.submodule {
      options = {
        enable = mkEnableOption "add a stats.yourdomain.tld";
        port = mkOption {
          type = types.number;
          default = 7890;
        };
      };
    };
  webOptions = with lib;
    types.submodule {
      options = {
        server = mkOption {type = serverOptions;};
        database = mkOption {type = databaseOptions;};
        goaccess = mkOption {type = goaccessOptions;};
        disableSSL = mkEnableOption "don't use ssl, usefull for internal testing";
      };
    };
in {
  options = {
    services.website = with lib; {
      enable = mkEnableOption "A website written in go :)";
      database = mkOption {
        description = "database root config";
        type = rootDatabaseOptions;
      };
      instances = mkOption {
        type = with types; attrsOf webOptions;
      };
    };
  };

  config = lib.mkIf cfg.enable {
    environment.systemPackages = [
      website
      pkgs.surrealdb
    ];
    services.surrealdb = {
      enable = true;
      port = 7999;
      extraFlags = [
        "-A" # enable all features
        "-u=${cfg.database.username}"
        "-p=${cfg.database.password}"
        "-l=full"
      ];
    };
    users.users.webhost = {
      isSystemUser = true;
      home = "/home/webhost";
      createHome = true;
      group = "webhost";
    };
    users.groups.webhost = {};
    services.caddy.enable = true;
    services.caddy.virtualHosts =
      lib.attrsets.concatMapAttrs (n: v: let
        ssl =
          if v.disableSSL
          then "http://"
          else "";
      in {
        "${ssl}${n}" = {
          serverAliases = [
            "${ssl}www.${n}"
          ];
          extraConfig = ''
            reverse_proxy "http://127.0.0.1:${toString v.server.port}"
          '';
        };
        "${ssl}stats.${n}" = lib.mkIf v.goaccess.enable {
          extraConfig = ''
            root * /var/lib/caddy/goaccess/${n}/
            file_server * browse

            @websockets {
              header Connection *Upgrade*
              header Upgrade websocket
            }
            reverse_proxy @websockets http://127.0.0.1:${toString v.goaccess.port}
          '';
        };
      })
      cfg.instances;
    systemd.services =
      lib.attrsets.concatMapAttrs (n: v: let
        safeN = lib.replaceStrings ["."] ["-"] n;
        dbname = "web";
        webconfig = (pkgs.formats.yaml {}).generate "config.yml" {
          server = {
            port = v.server.port;
            hostname = n;
          };
          database = {
            uri = let
              cfg = config.services.surrealdb;
            in "ws://${cfg.host}:${toString cfg.port}/rpc";
            namespace = n;
            name = dbname;
            username = safeN;
          };
        };
      in {
        "${safeN}-prep-db" = {...}: {
          description = "Setup database credentials";
          wantedBy = ["multi-user.target"];
          after = ["network.target" "surrealdb.service"];
          serviceConfig = {
            ExecStartPre = "${pkgs.coreutils-full}/bin/sleep 2"; # last-ditch effort to get db connecting first try
            ExecStart = let
              surreal = let
                cfgdb = config.services.surrealdb;
              in "${pkgs.surrealdb}/bin/surreal sql -e http://${cfgdb.host}:${toString cfgdb.port} -u ${cfg.database.username} -p ${cfg.database.password}";
              sql = ''
                DEFINE NAMESPACE IF NOT EXISTS ${n};
                USE NAMESPACE ${n};
                DEFINE DATABASE IF NOT EXISTS ${dbname};
                USE DATABASE ${n};
                DEFINE USER OVERWRITE ${safeN} ON NAMESPACE PASSWORD \"''${DB_PASSWORD}\" ROLES EDITOR;
              '';
              # verify with `echo "use ns localhost;use db localhost;info for ns;" | surreal sql -e ws://127.0.0.1:7999 -u root -p root`
              start = pkgs.writeShellScriptBin "start" ''
                echo "${sql}" | ${surreal}
              '';
            in "${start}/bin/start";
            DynamicUser = true;
            Type = "oneshot";
            Restart = "no";
          };
          environment = {
            DB_PASSWORD = v.database.password;
          };
        };
        "${safeN}-site" = {...}: {
          description = "A website written in go :";
          wantedBy = ["multi-user.target"];
          after = ["network.target" "${safeN}-prep-db.service"];
          serviceConfig = {
            ExecStart = "${website}/bin/web -static ${website} -state /var/lib/${safeN} -config ${webconfig}";
            DynamicUser = false;
            User = "webhost";
            Group = "webhost";
            StateDirectory = n;
            StateDirectoryMode = "0750";
          };
          environment = {
            DB_PASSWORD = v.database.password;
          };
        };
        "${safeN}-goaccess" = lib.mkIf v.goaccess.enable ({...}: {
          description = "an open source real-time web log analyzer";
          wantedBy = ["multi-user.target"];
          after = ["network.target"];
          environment = {
          };
          serviceConfig = {
            ExecStart = let
              mkdir = "${pkgs.coreutils-full}/bin/mkdir";
              touch = "${pkgs.coreutils-full}/bin/touch";
              accessLog = "/var/log/caddy/access-${n}.log";
              hostDir = "/var/lib/caddy/goaccess/${n}";
              hostFile = "${hostDir}/index.html";
              wsUrl =
                if v.disableSSL
                then "ws://stats.${n}:80"
                else "wss://stats.${n}:443";
              goaccess = pkgs.writeShellScriptBin "start" ''
                ${mkdir} -p ${hostDir}
                ${touch} ${hostFile}
                ${pkgs.goaccess}/bin/goaccess --log-format caddy -f ${accessLog} -o ${hostFile} --real-time-html --port=${toString v.goaccess.port} --ws-url=${wsUrl} --anonymize-ip --anonymize-level=3
              '';
            in "${goaccess}/bin/start";
            DynamicUser = false;
            User = "caddy";
            Group = "caddy";
          };
        });
      })
      cfg.instances;
  };
}
