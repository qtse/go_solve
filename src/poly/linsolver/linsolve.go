package linsolver

import (
    "fmt"
    "poly"
    )

type Solution map[string]float64

func (s Solution) String() string {
  res := ""
  first := true
  for p,v := range s {
    if !first {
      res += ", "
    } else {
      first = false
    }

    res += fmt.Sprintf("%s:%f",p,v)
  }
  return res
}

func Solve(system []poly.Polynomial) Solution {
  soln := make(map[string]float64)

  for _,p := range system {
    for _,t := range p {
      for pron,_ := range t.Pron {
        soln[pron] = 0
      }
    }
  }
  return soln
}

type Polynomial poly.Polynomial
