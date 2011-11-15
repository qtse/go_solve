package main

import (
    "poly"
    "poly/linsolver"
    )

func main() {
  defer cleanUp()
  p := poly.NewPronumeral("p")
  q := poly.NewPronumeral("q")
  r := poly.NewPronumeral("r")
  one := poly.NewConstant(-1)
  two := poly.NewConstant(-2)
  system := make([]poly.Polynomial,3)
  system[0] = poly.AddPoly(p, q)
  system[1] = poly.MultPoly(p, two)
  system[1].AddPoly(q).AddPoly(one)
  system[2] = poly.AddPoly(r,two).AddPoly(p)
  println(system[0].String())
  println(system[1].String())
  println(system[2].String())
  res,err := linsolver.Solve(system)
  if err != nil {
    println(res.String())
    println(err.String())
  } else {
    println(res.String())
  }
}

func cleanUp() {
  linsolver.Done()
  poly.Done()
}
