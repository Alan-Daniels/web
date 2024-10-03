flake: {
  pkgs,
  system,
  lib,
  config,
  ...
}: let
  website = flake.packages."${system}".default;
  cfg = config.services.website;
  webOptions = with lib;
    types.submodule {
      options = {
        port = mkOption {
          type = types.number;
        };
        disableSSL = mkEnableOption "don't use ssl, usefull for internal testing";
        enablePublicStats = mkEnableOption "show stats on stats.yourdomain.tld";
        statsPort = mkOption {
          type = types.number;
          default = 7890;
        };
      };
    };
in {
  options = {
    services.website = with lib; {
      enable = mkEnableOption "A website written in go :)";
      instances = mkOption {
        type = with types; attrsOf webOptions;
      };
    };
  };

  config = lib.mkIf cfg.enable {
    environment.systemPackages = [
      website
    ];
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
        ssl = if v.disableSSL then "http://" else "";
      in {
        "${ssl}${n}" = {
          serverAliases = [
            "${ssl}www.${n}"
          ];
          extraConfig = ''
            reverse_proxy "http://127.0.0.1:${toString v.port}"
          '';
        };
        "${ssl}stats.${n}" = lib.mkIf v.enablePublicStats {
          extraConfig = ''
            root * /var/lib/caddy/goaccess/${n}/
            file_server * browse

            @websockets {
              header Connection *Upgrade*
              header Upgrade websocket
            }
            reverse_proxy @websockets http://127.0.0.1:${toString v.statsPort}
          '';
        };
      })
      cfg.instances;
    systemd.services =
      lib.attrsets.concatMapAttrs (n: v: let
        prettyN = lib.replaceStrings ["."] ["-"] n;
        logdir = "/home/webhost"; #todo: maybe move this?
      in {
        "${prettyN}-site" = {...}: {
          description = "A website written in go :";
          wantedBy = ["multi-user.target"];
          after = ["network.target"];
          serviceConfig = {
            ExecStart = let
              procs = pkgs.writeShellScriptBin "start" ''
                set -e
                ${website}/bin/web -root ${website} -logdir ${logdir} -port ${toString v.port}
              '';
            in "${procs}/bin/start";
            DynamicUser = false;
            User = "webhost";
            Group = "webhost";
          };
        };
        "${prettyN}-goaccess" = {...}: {
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
              wsUrl = if v.disableSSL then "ws://stats.${n}:80" else "wss://stats.${n}:443";
              goaccess = pkgs.writeShellScriptBin "start" ''
                ${mkdir} -p ${hostDir}
                ${touch} ${hostFile}
                ${pkgs.goaccess}/bin/goaccess --log-format caddy -f ${accessLog} -o ${hostFile} --real-time-html --port=${toString v.statsPort} --ws-url=${wsUrl} --anonymize-ip --anonymize-level=3
              '';
            in "${goaccess}/bin/start";
            DynamicUser = false;
            User = "caddy";
            Group = "caddy";
          };
        };
      })
      cfg.instances;
  };
}
