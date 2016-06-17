// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
)

// extids will be a flag receiver for adding chains and entries
type extids []string

func (e *extids) String() string {
	return fmt.Sprint(*e)
}

func (e *extids) Set(s string) error {
	*e = append(*e, s)
	return nil
}

var addchain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli mkchain [-e EXTID1 -e EXTID2 ...] ECADDRESS <STDIN>"
	cmd.description = "Create a new factom chain. Read data for the first entry from stdin. Use the entry credits from the specified address."
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			eids extids
		)
		flag.Var(&eids, "e", "external id for the entry")
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		e := new(factom.Entry)

		for _, id := range eids {
			e.ExtIDs = append(e.ExtIDs, []byte(id))
		}

		// Entry.Content is read from stdin
		if p, err := ioutil.ReadAll(os.Stdin); err != nil {
			errorln(err)
			return
		} else if size := len(p); size > 10240 {
			errorln(fmt.Errorf("Entry of %d bytes is too large", size))
			return
		} else {
			e.Content = p
		}

		c := factom.NewChain(e)

		if _, err := factom.GetChainHead(c.ChainID); err == nil {
			// no error means the client found the chain
			errorln("Chain", c.ChainID, "already exists")
			return
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}
		// commit the chain
		if err := factom.CommitChain(c, ec); err != nil {
			errorln(err)
			return
		}
		// get commit acknowledgement
		// reveal chain
		// ? get reveal ack
	}
	help.Add("addchain", addchain)
	return cmd
}()