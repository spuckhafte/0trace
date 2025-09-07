{
  config,
  pkgs,
  lib,
  ...
}: {
  # Base system for ISO (you can pick text-based or graphical)
  imports = [<nixpkgs/nixos/modules/installer/cd-dvd/installation-cd-minimal.nix>];

  # Force autologin as root
  services.getty.autologinUser = lib.mkForce "root";

  # Enable NetworkManager
  networking.networkmanager.enable = true;
  networking.useDHCP = false;

  # If you want to include your custom TUI binary via overlay
  nixpkgs.overlays = [(import ./overlay.nix)];

  # Make sure your binary and other tools are included
  environment.systemPackages = with pkgs; [
    ztrace # from your overlay UPDATEME
    bash
    coreutils
  ];

  # Run your TUI automatically at startup
  systemd.user.services.ztrace = {
    description = "Auto-run TUI App";
    wantedBy = ["default.target"];
    serviceConfig = {
      ExecStart = "${pkgs.ztrace}/bin/mytui";
      StandardInput = "tty";
      StandardOutput = "tty";
      TTYPath = "/dev/tty1";
      Restart = "always";
    };
  };
}
