package linsolver

import (
    "fmt"
    "os"
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

func Solve(system []poly.Polynomial) (Solution, os.Error) {
  _,symTbl,err := formMatrix(system)
  if err != nil {
    return nil,err
  }
  soln := make(Solution,len(symTbl))
  todo := make([]int,len(symTbl)+1,len(symTbl)+1)
  for i,_ := range todo {
    todo[i] = i
  }

  for _,p := range system {
    for _,t := range p {
      for pron,_ := range t.Pron {
        soln[pron] = 0
      }
    }
  }
  return soln,os.NewError("Not implemented")
}

func formMatrix(system []poly.Polynomial) (m matrix, symTbl map[string]int, err os.Error) {
  symTbl = make(map[string]int)
  idx := 0
  for _,p := range system {
    for _,t := range p {
      for pron,_ := range t.Pron {
        if _,found := symTbl[pron]; !found {
          symTbl[pron]=idx
          idx++
        }
      }
    }
  }

  m = make(matrix,len(system),len(system))
  for i,p := range system {
    m[i] = make(row,len(symTbl)+1,len(symTbl)+1)

    for _,t := range p {
      switch t.Degree() {
      case 0:
        idx = len(symTbl)
      case 1:
        for pron,_ := range t.Pron {
          // This should only be executed once
          var ok bool
          if idx,ok = symTbl[pron]; !ok {
            err = os.NewError("internal error - pronumeral not found")
            return
          }
        }
      default:
        err = os.NewError("Non linear equation")
        return
      }
      m[i][idx] = t.Coeff
    }
  }
  return
}

func (r row) subtract(ref row, coeff float64) (res row, err os.Error) {
  if len(r) != len(ref) {
    return nil,os.NewError("row length mismatch")
  }
  for i := range r {
    r[i] -= coeff * ref[i]
  }
  return r,nil
}

func (r row) normalise(scale float64) row {
  for i := range r {
    r[i] *= scale
  }
  return r
}

func (m matrix) GetSolution(symTbl map[string]int) (Solution,os.Error) {
  soln := make(Solution,len(symTbl))

  for _,r := range m {
    //TODO
    r.subtract(r,0)
  }
  return soln,nil
}

type row []float64
type matrix []row
