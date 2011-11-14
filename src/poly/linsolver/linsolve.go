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
  m,symTbl,err := formMatrix(system)
  if err != nil {
    return nil,err
  }
println(m.String())
  soln := make(Solution,len(symTbl))
  todo := make(map[string]int)
  for p,i := range symTbl {
    todo[p] = i
  }

  for ri,r := range m {
    if v,c,ok,z := r.solved(todo);ok {
      if !z {
        soln[v] = c
        todo[v] = -1,false
      }
      continue
    }
    vIdx := -1
    for i,c  := range r[:len(r)-1] {
      if c == 0 {
        continue
      } else {
        vIdx = i
        break
      }
    }
    if vIdx < 0 {
      // Useless row
      continue
    }
    r.normalise(1/r[vIdx])
    for ri2,r2 := range m {
      if ri2 == ri {
        continue
      }
      if v,c,ok,z := r2.solved(todo);ok {
        if !z {
          soln[v] = c
          todo[v] = -1,false
        }
        continue
      }
      r2.subtract(r,r2[vIdx])
      if v,c,ok,_ := r2.solved(todo);ok {
        soln[v] = c
        todo[v] = -1,false
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

func (r row) solved(todoTbl map[string]int) (string,float64,bool,bool) {
  pStr := ""
  for v,i := range todoTbl {
    if pStr != "" && r[i] != 0 {
      return "",0,false,false
    } else if r[i] != 0 {
      pStr = v
    }
  }
  if pStr == "" {
    return pStr,r[len(r)-1],true,true
  }
  return pStr,r[len(r)-1],true,false
}

func (m matrix) String() (res string) {
  for _,r := range m {
    for _,c := range r {
      res += fmt.Sprintf("%f\t",c)
    }
    res += fmt.Sprintln()
  }
  return
}

type row []float64
type matrix []row
