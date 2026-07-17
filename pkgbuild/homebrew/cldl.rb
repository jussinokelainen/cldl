class Cldl < Formula
  desc "Directory based cli todo lists"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.3.0.tar.gz"
  sha256 "b6a604c919cf36a2acd9ebc1d9784b1f672a39fa12e438e0eb15d2997690710e"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make","PREFIX=#{prefix}", "install"
  end
end
