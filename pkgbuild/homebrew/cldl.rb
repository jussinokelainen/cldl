class Cldl < Formula
  desc "Directory based cli todo lists"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.2.4.tar.gz"
  sha256 "f4e3c22c5d25ee460436fc001e0772915cf9f8594c49c779a7b79f9856d7c9e3"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make", "release"
    bin.install "build/cldl"
    pkgshare.install "default_config.toml"
  end
end
