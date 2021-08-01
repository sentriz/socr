# Changelog

## [0.3.0](https://www.github.com/sentriz/socr/compare/v0.2.2...v0.3.0) (2021-08-01)


### Features

* add db migrations ([86c2db1](https://www.github.com/sentriz/socr/commit/86c2db17007a06dd0635d252c971dbb78989c061))
* add filter by media to search endpoint ([4dede23](https://www.github.com/sentriz/socr/commit/4dede234567ce07bd5fbb7bbffad030f0aad13b0))
* add image preview thumbnails ([980c09a](https://www.github.com/sentriz/socr/commit/980c09a44297ca9f89f54cd0e9950c04678f0a23))
* **ci:** check prettier ([26f1803](https://www.github.com/sentriz/socr/commit/26f18035e7a68439f6090256a340c9c02dd6c6e0))
* **ci:** lint and test frontend ([8435204](https://www.github.com/sentriz/socr/commit/84352040043152fcce4da6bd73ed0c97effd4e24))
* **readme:** add flow img ([d5757d8](https://www.github.com/sentriz/socr/commit/d5757d8e990a542dbd1f01b1fe2e29cec184f127))
* **ui:** add "check" command ([3878c6c](https://www.github.com/sentriz/socr/commit/3878c6cc3253932b4de83b161ec281308026d0ad))
* **ui:** add no results found state ([f9f1652](https://www.github.com/sentriz/socr/commit/f9f16520b82a50bd843885e2f44703ac496b5e04))
* **ui:** add option to filter by media type ([01ff10a](https://www.github.com/sentriz/socr/commit/01ff10a74a210ad92d1f80356e8d709472eab76e))
* **ui:** don't show public text for videos ([12a036e](https://www.github.com/sentriz/socr/commit/12a036e847dd137fd89490d38007a2aa81f563c3))
* **ui:** show video icon for videos ([ad03b6d](https://www.github.com/sentriz/socr/commit/ad03b6dedf8bcddd7991731b8cb1760c831fa22f))
* **ui:** use hericons ([948f459](https://www.github.com/sentriz/socr/commit/948f45918a3d857a05a23d7f58d13dfb3c3c6e53))


### Bug Fixes

* **build:** add build deps to final container ([a02943e](https://www.github.com/sentriz/socr/commit/a02943e63ef2014a195922739fc8a202307fc859))
* **lint:** remove funlen ([b102bea](https://www.github.com/sentriz/socr/commit/b102beab72d4e83f26328926d6870eeca8ea55ad))
* swap image / video icons ([75de539](https://www.github.com/sentriz/socr/commit/75de539d54bba53abc9a7b8510be4715cccc5158))
* **ui:** emit load event for video with loadstart ([ea5c2a0](https://www.github.com/sentriz/socr/commit/ea5c2a042086dc843a710108a01d218c33ef457b))
* **ui:** fix type errors ([25dfd0e](https://www.github.com/sentriz/socr/commit/25dfd0e8bdf2cd58e0e3eb779a12b2ef01e7f5d2))
* update accepted sort order dec -> desc ([cd53fe4](https://www.github.com/sentriz/socr/commit/cd53fe48922bfad041ba3e52ba35bb74f0c95467))

### [0.2.2](https://www.github.com/sentriz/socr/compare/v0.2.1...v0.2.2) (2021-05-14)


### Bug Fixes

* **ci:** install build deps for prod image ([f3528ba](https://www.github.com/sentriz/socr/commit/f3528ba688f485d68f6c494c0775d1a964e47198))

### [0.2.1](https://www.github.com/sentriz/socr/compare/v0.1.3...v0.2.1) (2021-05-14)


### Features

* add support for video ([7dfd0d8](https://www.github.com/sentriz/socr/commit/7dfd0d87eccb3dc50117425923846335160c6741))
* **ci:** pin golangci-lint version ([09be034](https://www.github.com/sentriz/socr/commit/09be03430647724ce15031ea371d4f031d804dbb))
* **ci:** test before release please, and only run extra tests on develop and pull request ([212587c](https://www.github.com/sentriz/socr/commit/212587c5348812d8f4413f4db12fcbc453c50712))
* **ci:** use GITHUB_TOKEN for release please ([57cbdd3](https://www.github.com/sentriz/socr/commit/57cbdd300c3b0f103a5481a0a337942bd65d8e04))
* **deps:** bump deps ([7614aeb](https://www.github.com/sentriz/socr/commit/7614aebee7e669000b008a1638f286a5f0fd8606))
* store mime and render video ([f8850b4](https://www.github.com/sentriz/socr/commit/f8850b45bc733fdf723755bf4b9a8e93aa3b8485))
* **ui:** reuse screenshot hightlight for public page ([988e2b1](https://www.github.com/sentriz/socr/commit/988e2b16f719264ec454a60968da2305be1c8b9f))
* **ui:** update frontend to use new terms and endpoints ([f9aa5a3](https://www.github.com/sentriz/socr/commit/f9aa5a3ba669853bae8093e62772471d21fb86f9))


### Bug Fixes

* **ci:** install ffmpeg deps ([7459cb3](https://www.github.com/sentriz/socr/commit/7459cb34b5281fe43f16c4699c1d72f75aac39bc))
* **ci:** trim short hash ([d3ade36](https://www.github.com/sentriz/socr/commit/d3ade36a62c34e00ad0f1ac610f912797eb8d7ff))
* **scanner:** try RFC3339 and add some shit tests ([a903449](https://www.github.com/sentriz/socr/commit/a903449c23ec7e918a0c0d09fb45e54280709452))

## [0.2.0](https://www.github.com/sentriz/socr/compare/v0.1.3...v0.2.0) (2021-05-12)


### Features

* **ci:** pin golangci-lint version ([09be034](https://www.github.com/sentriz/socr/commit/09be03430647724ce15031ea371d4f031d804dbb))
* **ci:** test before release please, and only run extra tests on develop and pull request ([212587c](https://www.github.com/sentriz/socr/commit/212587c5348812d8f4413f4db12fcbc453c50712))
* **ci:** use GITHUB_TOKEN for release please ([57cbdd3](https://www.github.com/sentriz/socr/commit/57cbdd300c3b0f103a5481a0a337942bd65d8e04))
* **deps:** bump deps ([7614aeb](https://www.github.com/sentriz/socr/commit/7614aebee7e669000b008a1638f286a5f0fd8606))


### Bug Fixes

* **ci:** trim short hash ([d3ade36](https://www.github.com/sentriz/socr/commit/d3ade36a62c34e00ad0f1ac610f912797eb8d7ff))

### [0.1.3](https://www.github.com/sentriz/socr/compare/v0.1.2...v0.1.3) (2021-05-08)


### Bug Fixes

* consistent release yaml ([0a8a2e9](https://www.github.com/sentriz/socr/commit/0a8a2e9e270589e3557c073c6a7e50c7854e9050))

### [0.1.2](https://www.github.com/sentriz/socr/compare/v0.1.1...v0.1.2) (2021-05-08)


### Bug Fixes

* show version on startup ([9eccd70](https://www.github.com/sentriz/socr/commit/9eccd70554aef1f3a1e5bacffdc191651d16ae5e))

### [0.1.1](https://www.github.com/sentriz/socr/compare/v0.1.0...v0.1.1) (2021-05-08)


### Bug Fixes

* **ci:** don't build arm v6 ([23835bc](https://www.github.com/sentriz/socr/commit/23835bcc9ddbedec93d63c3812d07d0142d8b903))

## 0.1.0 (2021-05-08)


### Features

* **ci:** arm builds ([1356eec](https://www.github.com/sentriz/socr/commit/1356eec1578e0ec68da954198b11261c6b8f65ce))


### Bug Fixes

* **ci:** test ([bd6fed4](https://www.github.com/sentriz/socr/commit/bd6fed43f79095695be87aaa50c65c5be07985dc))
