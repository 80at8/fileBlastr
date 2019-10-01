package main

import (
	"strconv"
	"testing"
	"time"
)

var testStack fileHashMap

func LoadTestData() {
	testStack.init("./test")
	testStack.files = make(map[string]fileHashEntry)

	hashes := []string{
		"shUwXvAzuyBGYCpG2ITnsxvwXeYLnTulJgMZIn2Br7M=",
		"5xwAPR0JW6P-NvJH4fRCTGtk3B4QI2tUFyqzuf_HVgY=",
		"36gYkp7piLu_aAoVtsjVFlZaM-rXvUzEHxKvJ39xXp8=",
		"M0t8BX7-J30VpfxVqKQUo0FezzP4C3uO-EeSIpUO9vk=",
		"OAuMSRnOV5rihz6EyabS-AKv0AgJ23iEVz8lv5juIVY=",
		"9fvFCu_aopqES5gNgtENdB5WQoFNZYjvpKoToMUojyM=",
		"0CR1O6FOq6vhTNmzCugQvV0JMToLHw4Ya-XLHdKsN1I=",
		"HVsExWtOSt8YiHbEuRaFyhhGnYzAKTEd2ek4hbvCows=",
		"YFPDD8sUQQHeb7JxbkJU025W6P_Q04UCtVkMLl34XSc=",
		"0MUkiyngA6ijyQsnP9EavgXmRnP_UOOZdKEt1AywzUM=",
	}
	for x := 0; x < 10; x++ {
		var testHash fileHashEntry
		testHash.fileHash = hashes[x]
		testHash.fileStatus = ACTIVE
		testHash.fileName = "testfile" + strconv.Itoa(x) + ".txt"
		testStack.add(testHash)
	}
}

func TestAdd(t *testing.T) {

	LoadTestData()

	if len(testStack.files) != 10 {
		t.Errorf("TestAdd(): add function did not add hashes, expected length of 10 found length of %d", len(testStack.files))
	}
}

func TestRemove(t *testing.T) {

	LoadTestData()

	for k := range testStack.files {
		delete(testStack.files, k)
	}

	if len(testStack.files) != 0 {
		t.Errorf("TestRemove(): did not remove all elements of array, expected 0 found %d", len(testStack.files))
	}
}

func TestRefresh(t *testing.T) {
	var x fileHashEntry
	LoadTestData()

	// adjust time to be in the past by 29 minutes. Files should be marked ACTIVE
	for k, v := range testStack.files {
		x = v
		x.entryAge = time.Now().Add(-29 * time.Minute)
		testStack.files[k] = x
	}
	testStack.refresh()

	n := 0
	for _, v := range testStack.files {
		if v.fileStatus != ACTIVE {
			n++
		}
	}

	if n != 0 {
		t.Errorf("TestRefresh(): Length of testStack.files is %d, Expected 0 STALE records, actual STALE records %d", len(testStack.files), n)
	}

	// adjust time to be in the past by 31 minutes. Files should be marked STALE
	for k, v := range testStack.files {
		x = v
		x.entryAge = time.Now().Add(-31 * time.Minute)
		testStack.files[k] = x
	}

	testStack.refresh()

	n = 0
	for _, v := range testStack.files {
		if v.fileStatus != STALE {
			n++
		}
	}

	if n != 0 {
		t.Errorf("TestRefresh(): Length of testStack.files is %d, Expected 0 ACTIVE records, actual ACTIVE records %d", len(testStack.files), n)
	}

}

func TestHash(t *testing.T) {
	LoadTestData()

	for _, v := range testStack.files {
		x := testStack.hash(v)
		_, v := testStack.files[x.fileHash]
		if !v {
			t.Errorf("TestHash(): Expected to find hash for %v but none found.", x.fileHash)
		}
	}

}
