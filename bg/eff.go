package bg

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"strconv"
	"unicode/utf8"
)

type effHeader struct {
	Signature, Version [4]byte
}

type EffEffect struct {
	Signature, Version [4]byte
	EffectID           uint32
	TargetType         uint32
	SpellLevel         uint32
	EffectAmount       int32
	DWFlags            uint32
	DurationType       uint32
	Duration           uint32
	ProbabilityUpper   uint16
	ProbabilityLower   uint16
	Res                [8]byte
	NumDice            uint32
	DiceSize           uint32
	SavingThrow        uint32
	SaveMod            int32
	Special            uint32

	School   uint32
	Unknown  uint32
	MinLevel uint32
	MaxLevel uint32
	Flags    uint32

	EffectAmount2 int32
	EffectAmount3 int32
	EffectAmount4 int32
	EffectAmount5 int32

	Res2 [8]byte
	Res3 [8]byte

	SourceX        int32
	SourceY        int32
	TargetX        int32
	TargetY        int32
	SourceType     uint32
	SourceRes      [8]byte
	SourceFlags    uint32
	ProjectileType uint32
	SlotNum        int32
	ScriptName     [32]byte
	CasterLevel    uint32
	FirstCall      uint32
	SecondaryType  uint32
	Pad            [15]uint32
}

func OpenEff(r io.ReadSeeker) (*ItmEffect, *EffEffect, error) {
	effHeader := &effHeader{}
	effv1 := &ItmEffect{}
	effv2 := &EffEffect{}

	err := binary.Read(r, binary.LittleEndian, effHeader)
	if err != nil {
		return nil, nil, err
	}

	buf := make([]byte, 1)
	_ = utf8.EncodeRune(buf, rune(effHeader.Version[1]))
	version, err := strconv.Atoi(string(buf))
	if err != nil {
		return nil, nil, err
	}

	if version == 1 {
		_, err = r.Seek(0, 0)
		if err != nil {
			return nil, nil, err
		}
		err = binary.Read(r, binary.LittleEndian, effv1)
		if err != nil {
			return nil, nil, err
		}
	} else if version == 2 {
		err = binary.Read(r, binary.LittleEndian, effv2)
		if err != nil {
			return nil, nil, err
		}
	}
	return effv1, effv2, nil
}

func (eff *ItmEffect) WriteJson(w io.Writer) error {
	bytes, err := json.MarshalIndent(eff, "", "\t")
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

func (eff *EffEffect) WriteJson(w io.Writer) error {
	bytes, err := json.MarshalIndent(eff, "", "\t")
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}
