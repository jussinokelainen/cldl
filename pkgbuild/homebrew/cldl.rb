class Cldl < Formula
  desc "Directory based cli todo lists"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.2.6.tar.gz"
  sha256 "90d3d05c73d2bdab5d18c5dd687389deeffd205549b18ecabc41ab3bdfbf5b3d"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make","PREFIX=#{prefix}", "install"
  end
end
