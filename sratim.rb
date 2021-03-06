# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Sratim < Formula
  desc ""
  homepage ""
  version "0.1.1"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/YakirOren/sratim/releases/download/v0.1.1/sratim_0.1.1_Darwin_arm64.tar.gz"
      sha256 "ca672168a727d0cb26b2f413158535f0e45713c364c9aa4234231b37d408d274"

      def install
        bin.install "sratim"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/YakirOren/sratim/releases/download/v0.1.1/sratim_0.1.1_Darwin_x86_64.tar.gz"
      sha256 "47bbb6b86b7cdeedf4e55455058c0d7608e6d226f9f2cdd17526b36ce7c7b744"

      def install
        bin.install "sratim"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/YakirOren/sratim/releases/download/v0.1.1/sratim_0.1.1_Linux_arm64.tar.gz"
      sha256 "b3e3cbb2a08676448dace7ab188c47fe00617b56aad54cc33a89c72b35acb391"

      def install
        bin.install "sratim"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/YakirOren/sratim/releases/download/v0.1.1/sratim_0.1.1_Linux_x86_64.tar.gz"
      sha256 "c085ee570c26acb230b4daa6ceb00959fe3708ae5c6cbb0dbfdb6dd04bc090cc"

      def install
        bin.install "sratim"
      end
    end
  end
end
