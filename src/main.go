package main

import (
    "poly"
    "poly/linsolver"
    )

func main() {
  p := poly.NewPronumeral("p")
  q := poly.NewPronumeral("q")
  one := poly.NewConstant(-1)
  two := poly.NewConstant(-2)
  system := make([]poly.Polynomial,2)
  system[0] = poly.AddPoly(p, q)
  system[1] = poly.MultPoly(p, two)
  system[1].AddPoly(q).AddPoly(one)
  println(system[0].String())
  println(system[1].String())
  res,_ := linsolver.Solve(system)
///  if err != nil {
///  } else {
    println(res.String())
///  }
}
