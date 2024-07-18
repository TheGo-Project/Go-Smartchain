package core

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// MyCircuit defines the circuit for our ZKP
type MyCircuit struct {
	A frontend.Variable
	B frontend.Variable
	C frontend.Variable
}

// Define the constraints of the circuit
func (circuit *MyCircuit) Define(api frontend.API) error {
	sum := api.Add(circuit.A, circuit.B)
	api.AssertIsEqual(sum, circuit.C)
	return nil
}

// Setup creates the proving and verification keys
func Setup() (*groth16.ProvingKey, *groth16.VerifyingKey, error) {
	circuit := &MyCircuit{}
	// Use the correct curve instance
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return nil, nil, err
	}
	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		return nil, nil, err
	}
	return &pk, &vk, nil
}
