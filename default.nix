{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
      import (fetchTree nixpkgs.locked) {
        overlays = [
          (import "${fetchTree gomod2nix.locked}/overlay.nix")
        ];
      }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
  self,
}:
buildGoApplication {
  pname = "website";
  version = "0.1";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  ldflags = [
  ];
  preBuild = let
    rev =
      if (self ? rev)
      then self.rev
      else self.dirtyRev;
  in ''
    echo ${rev} > internal/commit.txt
    go run github.com/a-h/templ/cmd/templ generate
  '';
  postInstall = let
  in ''
    mkdir $out/assets
    cp -r ${./assets}/* $out/assets
  '';
}
