# Copyright 2026 Jeremy Edwards
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

presubmit: lint check

lint:
	go fmt ./...
	go vet ./...

test: check

check:
	CGO_ENABLED=0 go test ./...
	CGO_ENABLED=1 go test -race ./...

RM = rm
ZIP = zip
RAR = rar
TAR = tar
SEVENZIP = 7z

TEST_ASSETS = testing/testassets/files/index.html
TEST_ASSETS += testing/testassets/files/site.js
TEST_ASSETS += testing/testassets/files/weird\ \#1.txt
TEST_ASSETS += testing/testassets/files/weird\ \#.txt
TEST_ASSETS += testing/testassets/files/weird$$.txt
TEST_ASSETS += testing/testassets/files/assets/1.txt
TEST_ASSETS += testing/testassets/files/assets/2.txt
TEST_ASSETS += testing/testassets/files/assets/onetwothree/1.txt
TEST_ASSETS += testing/testassets/files/assets/onetwothree/2.txt
TEST_ASSETS += testing/testassets/files/assets/onetwothree/3.txt
TEST_ASSETS += testing/testassets/files/assets/four/4.txt
TEST_ASSETS += testing/testassets/files/assets/sixseven/6.txt
TEST_ASSETS += testing/testassets/files/assets/sixseven/7.txt
TEST_ASSETS += testing/testassets/files/assets/images/1.txt
TEST_ASSETS += testing/testassets/files/assets/images/2.txt
TEST_ASSETS += testing/testassets/files/assets/images/laptop.png
TEST_ASSETS += testing/testassets/files/assets/images/walking-duck.gif
TEST_ASSETS += testing/testassets/files/assets/images/eXample.TIFF
TEST_ASSETS += testing/testassets/files/assets/images/blue.ico
TEST_ASSETS += testing/testassets/files/assets/images/Laptop.jpg
TEST_ASSETS += testing/testassets/files/assets/images/earth.avi
TEST_ASSETS += testing/testassets/files/assets/images/earth.mp4
TEST_ASSETS += testing/testassets/files/assets/images/earth.webm


ARCHIVE_TEST_ASSETS = testing/testassets/archives/nodir-testassets.zip
ARCHIVE_TEST_ASSETS += testing/testassets/archives/single-testassets.zip
ARCHIVE_TEST_ASSETS += testing/testassets/archives/nested-testassets.zip
ARCHIVE_TEST_ASSETS += testing/testassets/archives/testassets.tar.gz
ARCHIVE_TEST_ASSETS += testing/testassets/archives/testassets.tar.bz2
ARCHIVE_TEST_ASSETS += testing/testassets/archives/testassets.tar.xz
ARCHIVE_TEST_ASSETS += testing/testassets/archives/testassets.tar.lz4
ARCHIVE_TEST_ASSETS += testing/testassets/archives/testassets.tar
ARCHIVE_TEST_ASSETS += testing/testassets/archives/testassets.7z

testassets: $(TEST_ASSETS)
archiveassets: $(ARCHIVE_TEST_ASSETS)
presubmit: lint check

testing/testassets/archives/nodir-testassets.zip: $(TEST_ASSETS)
	mkdir -p $(dir $@)
	cd testing/testassets/files/assets/onetwothree/; $(ZIP) -qr9 ../../../archives/nodir-testassets.zip 1.txt 2.txt 3.txt

testing/testassets/archives/single-testassets.zip: $(TEST_ASSETS)
	mkdir -p $(dir $@)
	cd testing/testassets/files/; $(ZIP) -qr9 ../archives/single-testassets.zip .

testing/testassets/archives/nested-testassets.zip: $(TEST_ASSETS) testing/testassets/archives/single-testassets.zip
	mkdir -p $(dir $@)
	cd testing/testassets/files/; $(ZIP) -qr9 ../archives/nested-testassets.zip .; $(ZIP) -qr9j ../archives/nested-testassets.zip ../archives/single-testassets.zip

testing/testassets/archives/testassets.tar.gz:
	mkdir -p $(dir $@)
	cd testing/testassets/files/; $(TAR) -I 'gzip -9' -cf ../archives/testassets.tar.gz *

testing/testassets/archives/testassets.tar.bz2:
	mkdir -p $(dir $@)
	cd testing/testassets/files/; BZIP=-9 $(TAR) cjf ../archives/testassets.tar.bz2 *

testing/testassets/archives/testassets.tar.xz:
	mkdir -p $(dir $@)
	cd testing/testassets/files/; $(TAR) cJf ../archives/testassets.tar.xz *

testing/testassets/archives/testassets.tar.lz4:
	mkdir -p $(dir $@)
	cd testing/testassets/files/; $(TAR) cf ../archives/testassets.tar.lz4 -I 'lz4' *

testing/testassets/archives/testassets.tar:
	mkdir -p $(dir $@)
	cd testing/testassets/files/; $(TAR) cf ../archives/testassets.tar *

testing/testassets/archives/testassets.7z:
	mkdir -p $(dir $@)
	cd testing/testassets/files/; $(SEVENZIP) a ../archives/testassets.7z *

testing/testassets/files/index.html:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/site.js:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/weird\ \#1.txt:
	mkdir -p testing/testassets/files/
	echo -n '$@' > '$@'

testing/testassets/files/weird\ \#.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/weird$$.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/1.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/2.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/onetwothree/1.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/onetwothree/2.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/onetwothree/3.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/four/4.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/sixseven/6.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/sixseven/7.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/images/1.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

testing/testassets/files/assets/images/2.txt:
	mkdir -p $(dir $@)
	echo -n $@ > $@

# https://picsum.photos/

testing/testassets/files/assets/images/laptop.png:
	mkdir -p $(dir $@)
	curl -L -o $@ "https://file-examples.com/storage/fe5938b8dd69855f49886e9/2017/10/file_example_PNG_500kB.png"

testing/testassets/files/assets/images/Laptop.jpg:
	mkdir -p $(dir $@)
	curl -L -o $@ "https://file-examples.com/storage/fe5938b8dd69855f49886e9/2017/10/file_example_JPG_100kB.jpg"

testing/testassets/files/assets/images/walking-duck.gif:
	mkdir -p $(dir $@)
	curl -L -o $@ "https://www.examplefile.com/images/downloaded.gif"

testing/testassets/files/assets/images/eXample.TIFF:
	mkdir -p $(dir $@)
	curl -L -o $@ "https://file-examples.com/storage/fe5938b8dd69855f49886e9/2017/10/file_example_TIFF_1MB.tiff"

testing/testassets/files/assets/images/blue.ico:
	mkdir -p $(dir $@)
	curl -L -o $@ "https://file-examples.com/storage/fe5938b8dd69855f49886e9/2017/10/file_example_favicon.ico"

testing/testassets/files/assets/images/earth.mp4:
	mkdir -p $(dir $@)
	curl -L -o $@ "https://file-examples.com/storage/fe5938b8dd69855f49886e9/2017/04/file_example_MP4_480_1_5MG.mp4"

testing/testassets/files/assets/images/earth.webm:
	mkdir -p $(dir $@)
	curl -L -o $@ "https://file-examples.com/storage/fe5938b8dd69855f49886e9/2020/03/file_example_WEBM_480_900KB.webm"

testing/testassets/files/assets/images/earth.avi:
	mkdir -p $(dir $@)
	curl -L -o $@ "https://file-examples.com/storage/fe5938b8dd69855f49886e9/2018/04/file_example_AVI_480_750kB.avi"

testing/testassets/files/assets/images/example.svg:
	mkdir -p $(dir $@)
	curl -L -o $@ "https://file-examples.com/wp-content/storage/2020/03/file_example_SVG_20kB.svg"

clean:
	rm -rf testing/testassets/

.PHONY: presubmit lint test check clean testassets
