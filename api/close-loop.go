package handler

import (
	"encoding/json"
	"math"
	"net/http"
)

// Constants and Data Structure
const eps = 1e-12

type Vec3 [3]float64

type Req struct {
	PTool6 Vec3 `json:"P_tool_6"`
	PToolF Vec3 `json:"P_tool_f"`
	S6F    Vec3 `json:"S6_f"`
	A67F   Vec3 `json:"a67_f"`
}

type Resp struct {
	A71     float64 `json:"a71"`
	S7      float64 `json:"S7"`
	S1      float64 `json:"S1"`
	Al71Deg float64 `json:"al71_deg"`
	Th7Deg  float64 `json:"th7_deg"`
	Gam1Deg float64 `json:"gaml_deg"`
}

// Helper Functions
func rad2deg(r float64) float64 {
	return r * 180 / math.Pi
}

func clamp1(x float64) float64 {
	if x < -1 {
		return -1
	} else if x > 1 {
		return 1
	} else {
		return x
	}
}

func dot(a, b Vec3) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func sub(a, b Vec3) (r Vec3) {
	r[0] = a[0] - b[0]
	r[1] = a[1] - b[1]
	r[2] = a[2] - b[2]
	return
}

func add(a, b Vec3) (r Vec3) {
	r[0] = a[0] + b[0]
	r[1] = a[1] + b[1]
	r[2] = a[2] + b[2]
	return
}

func scl(a Vec3, s float64) (r Vec3) {
	r[0] = a[0] * s
	r[1] = a[1] * s
	r[2] = a[2] * s
	return
}

func cross(a, b Vec3) (r Vec3) {
	r[0] = a[1]*b[2] - a[2]*b[1]
	r[1] = a[2]*b[0] - a[0]*b[2]
	r[2] = a[0]*b[1] - a[1]*b[0]
	return
}

func norm(a Vec3) float64 {
	return math.Sqrt(dot(a, a))
}

func normed(a Vec3) (r Vec3) {
	n := norm(a)
	if n < eps {
		return Vec3{0, 0, 0}
	}
	return Vec3{a[0] / n, a[1] / n, a[2] / n}
}

func projectPerp(a, n Vec3) (r Vec3) {
	n2 := dot(n, n)
	if n2 < eps {
		return Vec3{0, 0, 0}
	}
	k := dot(a, n) / n2
	return sub(a, scl(n, k))
}

func bearingAroundAxis(axis, nhat Vec3) float64 {
	X, Y := Vec3{1, 0, 0}, Vec3{0, 1, 0}
	xref := projectPerp(X, axis)
	if norm(xref) < eps {
		xref = projectPerp(Y, axis)
	}
	xref = normed(xref)
	yref := normed(cross(axis, xref))
	return math.Atan2(dot(yref, nhat), dot(xref, nhat))
}

// Close Loop Algorithm of the Puma Manipulator Inverse Kinematics
func closeLoop(req Req) Resp {
	P1 := Vec3{0, 0, 0}
	u := Vec3{0, 0, 1}
	P7 := req.PToolF
	v := normed(req.A67F)

	w := sub(P7, P1)
	uv := dot(u, v)
	al71 := rad2deg(math.Acos(clamp1(uv)))

	uxv := cross(u, v)

	if dot(uxv, uxv) < eps {
		S1 := dot(w, u)
		Pc1 := add(P1, scl(u, S1))
		nvec := sub(Pc1, P7)
		a71 := norm(nvec)
		S7 := dot(sub(P1, P7), v)
		th7, gam1 := 0.0, 0.0
		if a71 >= eps {
			nh := normed(nvec)
			th7 = rad2deg(bearingAroundAxis(v, nh))
			gam1 = rad2deg(bearingAroundAxis(u, nh))
		}
		return Resp{a71, S7, S1, al71, th7, gam1}
	}

	wu, wv := dot(w, u), dot(w, v)
	den := 1 - uv*uv
	s := (wu - wv*uv) / den
	t := (wu*uv - wv) / den
	C1 := add(P1, scl(u, s))
	C7 := add(P7, scl(v, t))
	nvec := sub(C1, C7)
	a71 := norm(nvec)

	th7, gam1 := 0.0, 0.0

	if a71 > eps {
		nh := normed(nvec)
		th7 = rad2deg(bearingAroundAxis(v, nh))
		gam1 = rad2deg(bearingAroundAxis(u, nh))
	}
	return Resp{a71, t, s, al71, th7, gam1}
}

// Site Entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp := closeLoop(req)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
