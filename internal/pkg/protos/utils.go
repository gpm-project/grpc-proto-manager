/**
 * Copyright 2023 GPM Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package protos

import (
	"os"
	"path"
	"path/filepath"

	"github.com/gpm-project/grpc-proto-manager/internal/pkg/files"
	"github.com/rs/zerolog/log"
)

// Common structure with functions used by the different implementations.
type Common struct {
}

// copyAllSourceFiles copies all source files into the generated path so it contains everything that will be uploaded
func (c *Common) copyAllSourceFiles(source string, generatedPath string) error {
	toCopy := make(map[string]string, 0)
	filepath.Walk(source, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			toCopy[currentPath] = info.Name()
		}

		return nil
	})
	for filePath, fileName := range toCopy {
		log.Debug().Str("toCopy", filePath).Str("fileName", fileName).Msg("moving file")
		err := files.CopyFile(filePath, path.Join(generatedPath, fileName))
		if err != nil {
			return err
		}
	}
	return nil
}

// moveGenerateFiles moves the generated files into the temp directory.
func (c *Common) moveGeneratedFiles(rootPath string, generatedPath string) error {
	log.Debug().Str("rootPath", rootPath).Str("generatedPath", generatedPath).Msg("moving generated content")

	// Find the generated files. This is a bit of a hack since the generated structure depends on language specs.
	// So we find all the files
	toCopy := make(map[string]string, 0)
	filepath.Walk(rootPath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			toCopy[currentPath] = info.Name()
		}

		return nil
	})

	for filePath, fileName := range toCopy {
		log.Debug().Str("toCopy", filePath).Str("fileName", fileName).Msg("moving file")

		err := c.moveFile(filePath, path.Join(generatedPath, fileName))
		if err != nil {
			return err
		}
	}

	// Cleanup the temporal generated directory.
	err := os.RemoveAll(rootPath)
	if err != nil {
		return err
	}

	return nil
}

// moveFile implements moving the contents of a file to a new path deleting the old one.
// os.Rename seems a reasonable alternative and works on Mac, however when testing
// on linux, the following error appeared cannot generate proto code: rename generated/github.com/dhiguero/grpc-internal-agenda-go/entities.pb.go /tmp/tmpFromConfig/grpc-internal-agenda-go/entities.pb.go: invalid cross-device link
func (c *Common) moveFile(from string, to string) error {
	// copy the file
	err := files.CopyFile(from, to)
	if err != nil {
		return err
	}
	// delete the old file
	err = os.RemoveAll(from)
	if err != nil {
		return err
	}
	return nil
}
