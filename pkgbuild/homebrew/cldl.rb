class Cldl < Formula
  desc "Directory based cli todo lists"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.2.5.tar.gz"
  sha256 "6b4de510bfe504338dfaf967c4412aa8b00ee7f4ebded4b12cdd0b3848e20473"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make", "release"
    bin.install "build/cldl"
    pkgshare.install "default_config.toml"
  end
end
