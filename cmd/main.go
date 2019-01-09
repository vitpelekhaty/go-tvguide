// IPTV guide viewer
//
// Copyright 2018 Vitaly Pelekhaty
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific
// language governing permissions and limitations under the License.

package main

import (
	"log"
	"os"
	"runtime"

	commands "go-tvguide/cmd/commands"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := commands.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
