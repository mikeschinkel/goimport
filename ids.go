package main

type Ids struct {
	dirId  int64
	repoId int64
}

func loadIdsFromInputs() (ids Ids) {
	ids = Ids{
		dirId:  selectDirFromLocalPath(inputs.rootPath),
		repoId: selectRepoIdFromOriginURL(inputs.gitOriginURL),
	}
	if ids.repoId == 0 {
		ids.repoId = insertRepoFromInputs()
	}
	if ids.dirId == 0 {
		ids.dirId = insertDirFromInputs(ids.repoId)
	}
	return ids
}
