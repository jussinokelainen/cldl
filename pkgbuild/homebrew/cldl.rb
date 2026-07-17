class Cldl < Formula
  desc "Directory based cli todo lists"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.2.7.tar.gz"
  sha256 "70b33be07bef98c92745f6aeb05b96dd7626730f82127d04fd8c9d4ae3c47989"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make","PREFIX=#{prefix}", "install"
  end
end
