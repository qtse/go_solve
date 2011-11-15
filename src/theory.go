package theory

import (
    "fmt"
    "math"
    "poly"
    )

func InitCalc(n int, cw int) Param {
  p := make([]poly.Polynomial,n+1,n+1)

  p[0] = nil // Unused
  for i := 1; i <= n; i++ {
    p[i] = poly.NewPronumeral(fmt.Sprintf("p%d",i))
  }

  return Param{n,cw,p}
}

func FormCondProbabilities(param Param) [][]poly.Polynomial {
  U := make([]float64,param.CW+1,param.CW+1)
  M := make([]float64,param.CW+1,param.CW+1)
  for i,_ := range U {
    U[i] = calcU(i,param.CW)
    M[i] = calcU(i,param.CW)
  }

  T := make([][]poly.Polynomial,param.CW + 1, param.CW + 1)

  for k := 0; k <= param.CW; k++ {
    T[k] = make([]poly.Polynomial,param.N+1,param.N+1)
    for will_tx := 0; will_tx <= param.N; will_tx++ {
      T[k][will_tx] = poly.NewConstant(0, true)
      for did_tx := 1; did_tx <= param.N; did_tx++ {
        did_not_tx := param.N - did_tx
        coeff := 0.

        for will_and_did_tx := 0;
                    will_and_did_tx <= did_tx;
                    will_and_did_tx++ {
          will_not_but_did_tx := did_tx - will_and_did_tx
          will_but_did_not_tx := will_tx - will_and_did_tx
          will_not_and_did_not_tx := did_not_tx - will_but_did_not_tx

          if will_but_did_not_tx < 0 {
            break
          } else if will_not_and_did_not_tx < 0 {
            continue
          }

          coeff += nchoosek(did_tx, will_and_did_tx) *
                      math.Pow(U[k], float64(will_and_did_tx)) *
                      math.Pow((1-U[k]), float64(will_not_but_did_tx)) *
                   nchoosek(did_not_tx, will_but_did_not_tx) *
                      math.Pow(M[k], float64(will_but_did_not_tx)) *
                      math.Pow((1-M[k]), float64(will_not_and_did_not_tx))
        }
        tmp := poly.MultConst(param.Polys[did_tx], coeff)
        T[k][will_tx] = poly.AddPoly(T[k][will_tx],tmp)
      }
    }
  }
  return T
}

func CalcApproxProb(T [][]poly.Polynomial, param Param) [][]poly.Polynomial {
  Tc := make([][]poly.Polynomial,param.CW+1,param.CW+1)

  for i := 0; i <= param.CW; i++ {
    Tc[i] = make([]poly.Polynomial,param.N+1,param.N+1)
    for j := 0; j <= param.N; j++ {
      if i == 0 {
        Tc[i][j] = T[i][j].Copy()
      } else {
        Tc[i][j] = poly.MultPoly(T[i][j], Tc[i-1][0])
      }
    }
  }

  return Tc
}

func SumProb(Tc [][]poly.Polynomial, param Param) []poly.Polynomial {
  P := make([]poly.Polynomial,param.N+1,param.N+1)
  P[0] = poly.NewConstant(0,true)
  for i := 1; i <= param.N; i++ {
    P[i] = poly.NewConstant(0,true).Copy()
    for j := 0; j <= param.CW; j++ {
      P[i].AddPoly(Tc[j][i])
    }
  }

  return P
}

func FormSystem(P []poly.Polynomial, param Param) []poly.Polynomial {
  system := make([]poly.Polynomial,param.N,param.N)
  system[0] = poly.NewConstant(-1,true).Copy()
  for _,p := range param.Polys[1:] {
    system[0].AddPoly(p)
  }

  for i := 1; i < param.N; i++ {
    system[1] = poly.AddPoly(P[i], poly.MultConst(param.Polys[i],-1))
  }

  return system
}

func calcU(k int, cw int) float64 {
  return 1/(float64(cw) + 1.0 - float64(k))
}

func calcM(k int, cw int) float64 {
  return 2/(float64(cw) + 2.0 - float64(k))
}

func nchoosek(n int, k int) float64 {
  return math.Gamma(float64(n+1))/(math.Gamma(float64(k+1))*math.Gamma(float64(n-k+1)))
}

type Param struct {
  N int
  CW int
  Polys []poly.Polynomial
}
