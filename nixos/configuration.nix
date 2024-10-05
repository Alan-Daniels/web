flake: {pkgs, ...}: let
in {
  imports = [
    ./hardware-configuration.nix
    flake.nixosModules.default
  ];

  nix.settings.experimental-features = ["nix-command" "flakes"];

  virtualisation.vmVariant = let
  in {
    # following configuration is added only when building VM with build-vm
    virtualisation = {
      memorySize = 512; # Use 0.5GB memory.
      cores = 1;
      diskImage = null; # don't save anything between boots
      forwardPorts = [
        {
          from = "host";
          host.port = 8080;
          guest.port = 8080;
        }
      ];
    };
    services.getty.autologinUser = "root";
    networking.firewall.allowedTCPPorts = [8080];

    virtualisation.qemu.options = [
      "-nographic"
    ];
  };

  time.timeZone = "Australia/Sydney";
  i18n.defaultLocale = "en_GB.UTF-8";
  i18n.extraLocaleSettings = {
    LC_ADDRESS = "en_AU.UTF-8";
    LC_IDENTIFICATION = "en_AU.UTF-8";
    LC_MEASUREMENT = "en_AU.UTF-8";
    LC_MONETARY = "en_AU.UTF-8";
    LC_NAME = "en_AU.UTF-8";
    LC_NUMERIC = "en_AU.UTF-8";
    LC_PAPER = "en_AU.UTF-8";
    LC_TELEPHONE = "en_AU.UTF-8";
    LC_TIME = "en_AU.UTF-8";
  };

  environment.systemPackages = with pkgs; [
  ];

  services.website = {
    enable = true;
    database = {
      username = "root";
      password = "root";
    };
    instances = {
      "localhost" = {
        disableSSL = true;
        server = {
          port = 8080;
        };
        database = {
          password = "testing-testing-hi";
        };
        goaccess.enable = false;
      };
    };
  };

  system.stateVersion = "23.11";
}
