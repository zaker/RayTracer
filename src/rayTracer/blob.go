package rayTracer

import (
	"math"
	"sort"
)

type Poly struct {
	a, b, c, fDistance, fDeltaFInvSquare float64
}

type polynoms []*Poly

func (p polynoms) Len() int      { return len(p) }
func (p polynoms) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p *Poly) abs() float64 {
	return math.Sqrt(p.a*p.a + p.b*p.b + p.c*p.c)

}

func (p polynoms) Less(i, j int) bool { return p[i].abs() < p[j].abs() }

type ZoneTab struct {
	fCoef            float64
	fDeltaFInvSquare float64
	fGamma           float64
	fBeta            float64
}

var zoneTab []ZoneTab = []ZoneTab{
	{10.0, 0, 0, 0},
	{5.0, 0, 0, 0},
	{3.33333, 0, 0, 0},
	{2.5, 0, 0, 0},
	{2.0, 0, 0, 0},
	{1.66667, 0, 0, 0},
	{1.42857, 0, 0, 0},
	{1.25, 0, 0, 0},
	{1.1111, 0, 0, 0},
	{1.0, 0, 0, 0}}

var zoneNumber int = len(zoneTab)

func InitBlobZones() {
	fLastGamma, fLastBeta := 0.0, 0.0
	fLastInvRSquare := 0.0
	for i := 0; i < zoneNumber-1; i++ {
		fInvRSquare := 1.0 / zoneTab[i+1].fCoef
		zoneTab[i].fDeltaFInvSquare = fInvRSquare - fLastInvRSquare
		// fGamma is the ramp between the entry point and the exit point.
		// We only store the difference compared to the previous zone
		// that way we can reconstruct the estimate more easily later..
		temp := (fLastInvRSquare - fInvRSquare) / (zoneTab[i].fCoef - zoneTab[i+1].fCoef)
		zoneTab[i].fGamma = temp - fLastGamma
		fLastGamma = temp

		// fBeta is the value of the line approaching the curve for dist = 0 (f = fGamma * x + fBeta)
		// similarly we only store the difference with the fBeta of the previous curve
		zoneTab[i].fBeta = fInvRSquare - fLastGamma*zoneTab[i+1].fCoef - fLastBeta
		fLastBeta = zoneTab[i].fBeta + fLastBeta

		fLastInvRSquare = fInvRSquare
	}
	// The last zone acts as a simple terminator 
	// (no need to evaluate the field there, because we know that it exceed
	// the equipotential value.. by design)
	zoneTab[zoneNumber-1].fGamma = 0.0
	zoneTab[zoneNumber-1].fBeta = 0.0

	print(zoneTab[0].fDeltaFInvSquare)

}

func (b *Blob) isBlobIntersected(r Ray, t float64) (bResult bool, ouT float64) {
	// Having a static structure helps performance more than two times !
	// It obviously wouldn't work if we were running in multiple threads..
	// But it helps considerably for now
	polynomMap := make(polynoms, 0)

	ouT = t

	rSquare := float64(b.Size * b.Size)
	rInvSquare := b.invSizeSquare
	maxEstimatedPotential := 0.0

	// outside of all the influence spheres, the potential is zero
	A, B, C := 0.0, 0.0, 0.0

	for i := range b.Centers {
		currentPoint := b.Centers[i]
		tmP := r.Start
		vDist := currentPoint.Sub(Vector{tmP.X, tmP.Y, tmP.Z})
		A = 1.0
		B = -2.0 * r.Dir.MulV(vDist)
		C = vDist.MulV(vDist)

		// Accelerate delta computation by keeping common computation outside of the loop
		BSquareOverFourMinusC := 0.25*B*B - C
		MinusBOverTwo := -0.5 * B
		ATimeInvSquare := A * rInvSquare
		BTimeInvSquare := B * rInvSquare
		CTimeInvSquare := C * rInvSquare

		// the current sphere, has N zones of influences
		// we go through each one of them, as long as we've detected
		// that the intersecting ray has hit them
		// Since all the influence zones of many spheres
		// are imbricated, we compute the influence of the current sphere
		// by computing the delta of the previous polygon
		// that way, even if we reorder the zones later by their distance
		// on the ray, we can still have our estimate of 
		// the potential function.
		// What is implicit here is that it only works because we've approximated
		// 1/dist^2 by a linear function of dist^2
		for j := 0; j < zoneNumber-1; j++ {
			// We compute the "delta" of the second degree equation for the current
			// spheric zone. If it's negative it means there is no intersection
			// of that spheric zone with the intersecting ray
			fDelta := BSquareOverFourMinusC + zoneTab[j].fCoef*rSquare
			if fDelta < 0.0 {
				// Zones go from bigger to smaller, so that if we don't hit the current one,
				// there is no chance we hit the smaller one
				// print("fDelta", zoneTab[j].fCoef)
				break
			}
			sqrtDelta := math.Sqrt(fDelta)
			t0 := MinusBOverTwo - sqrtDelta
			t1 := MinusBOverTwo + sqrtDelta

			// because we took the square root (a positive number), it's implicit that 
			// t0 is smaller than t1, so we know which is the entering point (into the current
			// sphere) and which is the exiting point.
			poly0 := Poly{zoneTab[j].fGamma * ATimeInvSquare,
				zoneTab[j].fGamma * BTimeInvSquare,
				zoneTab[j].fGamma*CTimeInvSquare + zoneTab[j].fBeta,
				t0,
				zoneTab[j].fDeltaFInvSquare}
			poly1 := Poly{-poly0.a, -poly0.b, -poly0.c,
				t1,
				-poly0.fDeltaFInvSquare}

			maxEstimatedPotential += zoneTab[j].fDeltaFInvSquare

			// just put them in the vector at the end
			// we'll sort all those point by distance later
			polynomMap = append(polynomMap, &poly0)
			polynomMap = append(polynomMap, &poly1)
		}
	}
	// if maxEstimatedPotential < 1.0 {
	// 	println("mEP", maxEstimatedPotential)
	// }

	if /*len(polynomMap) < 2 ||*/ maxEstimatedPotential < 1.0 {
		// println("p,m", len(polynomMap), maxEstimatedPotential)
		return
	}

	// sort the various entry/exit points per distance
	// by going from the smaller distance to the bigger
	// we can reconstruct the field approximately along the way
	sort.Sort(polynomMap)

	maxEstimatedPotential = 0.0

	for i, it := range polynomMap {
		// A * x2 + B * y + C, defines the condition under which the intersecting
		// ray intersects the equipotential surface. It works because we designed it that way
		// (refer to the article).

		A += it.a
		B += it.b
		C += it.c
		maxEstimatedPotential += it.fDeltaFInvSquare
		// println("pm", len(polynomMap), maxEstimatedPotential)
		// if maxEstimatedPotential < 1.0 {
		if i == 0 {
			// No chance that the potential will hit 1.0f in this zone, go to the next zone
			// just go to the next zone, we may have more luck
			// print(maxEstimatedPotential)
			continue
		}

		// fZoneStart := it.fDistance
		fZoneStart := polynomMap[i-1].fDistance

		fZoneEnd := polynomMap[i].fDistance

		// the current zone limits may be outside the ray start and the ray end
		// if that's the case just go to the next zone, we may have more luck
		if t > fZoneStart && 0.01 < fZoneEnd {
			// This is the exact resolution of the second degree
			// equation that we've built
			// of course after all the approximation we've done
			// we're not going to have the exact point on the iso surface
			// but we should be close enough to not see artifacts
			fDelta := B*B - 4.0*A*(C-1.0)
			if fDelta < 0.0 {
				continue
			}

			fInvA := (0.5 / A)
			fSqrtDelta := math.Sqrt(fDelta)

			t0 := fInvA * (-B - fSqrtDelta)
			t1 := fInvA * (-B + fSqrtDelta)

			if t0 > 0.01 && t0 >= fZoneStart && t0 < fZoneEnd && t0 <= t {
				ouT = t0
				bResult = true
			}

			if t1 > 0.01 && t1 >= fZoneStart && t1 < fZoneEnd && t1 <= t {
				ouT = t1
				bResult = true
			}
			// print("Hello")
			if bResult {
				// return true;
				// print("Hello")
				return
			}
		}
	}
	return
}

func (b *Blob) BlobInterpolation(pos Point) Vector {
	gradient := &Vector{0.0, 0.0, 0.0}

	fRSquare := float64(b.Size * b.Size)
	for i := range b.Centers {
		// This is the true formula of the gradient in the
		// potential field and not an estimation.
		// gradient = normal to the iso surface
		normal := pos.SubP(b.Centers[i])
		fDistSquare := normal.MulV(normal)
		if fDistSquare <= 0.001 {
			continue
		}
		fDistFour := fDistSquare * fDistSquare
		normal = normal.MulC(fRSquare / fDistFour)

		gradient.AddEq(normal)
	}
	return *gradient
}
