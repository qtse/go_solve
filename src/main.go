package main

import (
    "fmt"
    "os"
    "poly"
    "poly/linsolver"
    "strconv"
    "theory"
    )

func main() {
  defer cleanUp()

  if len(os.Args) != 3 {
    os.Exit(1)
  }

  n,err := strconv.Atoi(os.Args[1])
  if err != nil {
    os.Exit(1)
  }

  cw,err := strconv.Atoi(os.Args[2])
  if err != nil {
    os.Exit(1)
  }

  param := theory.InitCalc(n,cw)
  condProb := theory.FormCondProbabilities(param)

  poly.SetMaxDegree(1)
  poly.SetMinCoeff(1e-8)

  approxProb := theory.CalcApproxProb(condProb, param)
  sumProb := theory.SumProb(approxProb,param)
  approxProb = nil

  system := theory.FormSystem(sumProb, param)
  linSoln,err := linsolver.Solve(system)

  fmt.Println(linSoln)
  if err != nil {
    println(err.String())
  }

  stub(system)

///  p := poly.NewPronumeral("p")
///  q := poly.NewPronumeral("q")
///  r := poly.NewPronumeral("r")
///  one := poly.NewConstant(-1)
///  two := poly.NewConstant(-2)
///  system := make([]poly.Polynomial,3)
///  system[0] = poly.AddPoly(p, q)
///  system[1] = poly.MultPoly(p, two)
///  system[1].AddPoly(q).AddPoly(one)
///  system[2] = poly.AddPoly(r,two).AddPoly(p)
///  println(system[0].String())
///  println(system[1].String())
///  println(system[2].String())
///  res,err := linsolver.Solve(system)
///  if err != nil {
///    println(res.String())
///    println(err.String())
///  } else {
///    println(res.String())
///  }
}

func cleanUp() {
  linsolver.Done()
  poly.Done()
}

func stub(a ...interface{}){}
