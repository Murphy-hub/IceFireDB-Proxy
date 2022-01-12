/*
 *
 *  * Licensed to the Apache Software Foundation (ASF) under one or more
 *  * contributor license agreements.  See the NOTICE file distributed with
 *  * this work for additional information regarding copyright ownership.
 *  * The ASF licenses this file to You under the Apache License, Version 2.0
 *  * (the "License"); you may not use this file except in compliance with
 *  * the License.  You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package main

import (
	"log"

	"github.com/alicebob/miniredis/v2"
)

func main() {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	log.Println(s.Addr())
	log.Println(s.Lpush("l2", "a"))
	log.Println(s.Lpush("l2", "b"))
	log.Println(s.List("l2"))
	log.Println(s.Lpop("l2"))
	log.Println(s.Lpop("l2"))
}
