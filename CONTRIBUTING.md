Contributing Guidelines

Getting Started:

Fork, and clone the repository onto your local drive.
Find an open issue and assign it to yourself.
LOCAL || git status to make sure you are on master
LOCAL || git remote add upstream https://github.com/sfbrigade/sfsbook.git
LOCAL || git pull --rebase upstream master
LOCAL || git checkout -b feature/nameOfIssueOrFeature
GITHUB || create initial pull request with task list of issues
GITHUB || update issue status to "In Progress"
After making edits and committing changes to your local feature branch:

LOCAL || git pull --rebase upstream master
LOCAL || git push origin nameOfIssueOrFeature
GITHUB || make a pull request merging your local feature branch into the
GITHUB || fill out git - explaining within pull request comment field
Collaborator / Reviewer:

REVIEWABLE.IO || Provide pull request URL to open review
REVIEWABLE.IO || Review and approve pull request
