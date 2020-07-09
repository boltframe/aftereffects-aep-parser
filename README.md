# aftereffects-aep-parser

This project is dedicated to reverse-engineering, parsing and exposing useful APIs for understanding Adobe After Effects AEP project files. This may be useful for a variety of reasons. For example, one might want to programmatically determine the number of compositions in a project, their resolution(s), framerate, etc. Traditionally, one would need to piggyback the ExtendScript engine running within After Effects to retrieve this information. While this procedure is fairly easy to accomplish, it is exponentially more resource expensive than static file analysis, and requires a running instance of the After Effects Scripting Engine (which can only be run on Windows and MacOS). 

# Quick Start

```bash
go get -u github.com/boltframe/aftereffects-aep-parser
```

```go
package main

import (
  "fmt"
  aep "github.com/boltframe/aftereffects-aep-parser"
)

func main() {
  project, err := aep.Open("./my-project.aep")
  if err != nil {
    panic(err)
  }
  fmt.Println(project)
}
```

# Contributing

Any and all contributions are welcome! There is no official procedure for contributing yet, but please start by opening a new issue describing your findings and anticipated contribution. Thank you!

# Research Procedure

After Effects Project files are currently encoded using the Resource Interchange File Format (RIFF) in a Big-Endian byte-ordering (also known as RIFX). This can be confirmed by inspecting the first four bytes of any AEP file which will contain the following ASCII characters: `RIFX....`. Once parsed into [Chunks](https://en.wikipedia.org/wiki/Resource_Interchange_File_Format#Explanation), you can inspect the internal representation and *attempt* to understand its structure.

I have found that the online tool [Kaitai](https://ide.kaitai.io) is extremely helpful for quickly viewing structured binary data. I use the following (possibly outdated) language definition for validating AEP files:

```ksy
meta:
  id: aep
  endian: be
  file-extension: aep

seq:
  - id: magic1
    contents: RIFX
  - id: file_size
    type: u4
  - id: magic2
    contents: Egg!
  - id: data
    type: blocks
    size: file_size - 4
    
types:
  blocks:
    seq:
      - id: entries
        type: block
        repeat: eos
  block:
    seq: 
      - id: block_type
        type: u4
        enum: chunk_type
      - id: block_size
        type: u4
      - id: data
        size: block_size
        type: 
          switch-on: block_type
          cases:
            'chunk_type::list': list_body
            'chunk_type::utf8': utf8_body
            'chunk_type::cdta': cdta_body
            'chunk_type::idta': idta_body
            'chunk_type::cmta': utf8_body
            'chunk_type::fdta': fdta_body
            _: ascii_body
      - id: padding
        type: u1
        if: (block_size % 2) != 0
  list_body:
    seq:
      - id: identifier
        type: str
        encoding: ascii
        size: 4
      - id: entries
        type: blocks
  utf8_body:
    seq:
      - id: data
        type: str
        encoding: utf8
        size-eos: true
  ascii_body:
    seq:
      - id: data
        type: str
        encoding: ascii
        size-eos: true
  idta_body:
    seq:
      - id: unknown1
        type: str
        size: 18
        encoding: ascii
      - id: id
        type: u2
  fdta_body:
    seq:
      - id: unknown1
        type: str
        encoding: ascii
        size: 1
  cdta_body:
    seq:
      - id: unknown1
        type: str
        size: 140
        encoding: ascii
      - id: width
        type: u2
      - id: height
        type: u2
      - id: unknown2
        type: str
        size: 12
        encoding: ascii
      - id: frame_rate
        type: u2
        
enums:
  chunk_type:
    0x4c495354: list
    0x55746638: utf8
    0x63647461: cdta # Composition data
    0x69647461: idta # Item data
    0x636d7461: cmta # Comment data
    0x66647461: fdta # Folder data
```

By making small changes to an After Effects project and uploading it to the Kaitai viewer, you can narrow down on how data is represented, and start to write additional language definitions and structures. 