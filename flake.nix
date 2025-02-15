{
  description = "Hey, what if Airtable, like... wasn't?";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
      rec {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            etcd
            postgresql
          ];
        };
        packages.landmine = pkgs.buildGoModule
          {
            pname = "landmine";
            version = "0.0.0";
            src = ./.;
            vendorHash = "sha256-mQ8aJbVXDAMYK76WDzC9HbcbzKMK07Wcaoo+VqhRGFA=";
          };
        packages.default = packages.landmine;
      }
    );
}
