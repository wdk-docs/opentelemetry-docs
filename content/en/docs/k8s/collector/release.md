# 收集器发布

收集器的构建和测试目前是完全自动化的。然而，为了发布，仍然需要手动执行某些操作。

我们以相同的版本发布 core 和 contrib 收集器，其中 contrib 版本使用 core 版本作为
依赖。我们把这个过程分为四个部分。发布工程师必须发布:

1. [Core](#releasing-opentelemetry-collector)收集器，包括收集器构建器 CLI 工具。
2. [Contrib](#releasing-opentelemetry-collector-contrib)收集器。
3. [artifacts](#producing-the-artifacts)

!!! danger

      你需要能够签署 git 提交/标签，以便能够发布收集器版本。
      按照[本指南](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)进行设置。

!!! danger

      为了能够发布，您需要成为两个存储库的审批人。
      这是必需的，因为您需要将标签直接推送并提交到以下存储库:

- [open-telemetry/opentelemetry-collector](https://github.com/open-telemetry/opentelemetry-collector)
- [open-telemetry/opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib)
- [open-telemetry/opentelemetry-collector-releases](https://github.com/open-telemetry/opentelemetry-collector-releases)

## 发布管理

发布管理员是负责特定发布的人。虽然管理员可能会向其他人请求帮助，但他们最终要对发
布的成功负责。

In order to have more people comfortable with the release process, and in order
to decrease the burden on a small number of volunteers, all core approvers are
release managers from time to time, listed under the
[Release Schedule](#release-schedule) section. That table is updated at every
release, with the current manager adding themselves to the bottom of the table,
removing themselves from the top of the table.

It is possible that a core approver isn't a contrib approver. In that case, the
release manager should coordinate with a contrib approver for the steps
requiring such role, like the publishing of tags.

## 发布 opentelemetry-collector

1. Determine the version number that will be assigned to the release. During the
   beta phase, we increment the minor version number and set the patch number
   to 0. In this document, we are using `v0.55.0` as the version to be released,
   following `v0.54.0`. Check if stable modules have any changes since the last
   release by running
   `make check-changes PREVIOUS_VERSION=v1.0.0-rc9 MODSET=stable`. If there are
   no changes, there is no need to release new version for stable modules.

2. Manually run the action
   [Automation - Prepare Release](https://github.com/open-telemetry/opentelemetry-collector/actions/workflows/prepare-release.yml).
   When prompted, enter the version numbers determined in Step 1, but do not
   include `v`. This action will create an issue to track the progress of the
   release and a pull request to update the changelog and version numbers in the
   repo. **While this PR is open all merging in Core should be haulted**. **Do
   not** move forward until this PR is merged.

   - If the PR needs updated in any way you can make the changes in a fork and
     PR those changes into the `prepare-release-prs/x` branch. You do not need
     to wait for the CI to pass in this prep-to-prep PR.

3. Update Contrib to use the latest in development version of Core. Run
   `make update-otel` in Contrib root directory and if it results in any changes
   submit a PR to Contrib with the changes as draft. This is to ensure that the
   latest core does not break contrib in any way. We’ll update it once more to
   the final release number later.

4. Create a branch named `release/<release-series>` (e.g. `release/v0.55.x`)
   from the `Prepare Release` commit and push to
   `open-telemetry/opentelemetry-collector`.

5. Make sure you are on `release/<release-series>`. Tag the module groups with
   the new release version by running:

   - `make push-tags MODSET=beta` for beta modules group,
   - `make push-tags MODSET=stable` beta stable modules group, only if there
     were changes since the last release.

   If you set your remote using `https` you need to include
   `REMOTE=https://github.com/open-telemetry/opentelemetry-collector.git` in
   each command. Wait for the new tag build to pass successfully.

6. The release script for the collector builder should create a new GitHub
   release for the builder. This is a separate release from the core, but we
   might join them in the future if it makes sense.

7. A new `v0.55.0` release should be automatically created on Github by now.
   Edit it and use the contents from the CHANGELOG.md as the release's
   description.

8. If you created a draft PR to Contrib in step 3, update the PR to use the
   newly released Core version and set it to Ready for Review. **Do not** move
   forward until this PR is merged.

## 发布 opentelemetry-collector-contrib

1. Manually run the action
   [Automation - Prepare Release](https://github.com/open-telemetry/opentelemetry-collector-contrib/actions/workflows/prepare-release.yml).
   When prompted, enter the version numbers determined in Step 1, but do not
   include `v`. This action will a pull request to update the changelog and
   version numbers in the repo. **While this PR is open all merging in Contrib
   should be haulted**. Do not move forward until this PR is merged.

   - If the PR needs updated in any way you can make the changes in a fork and
     PR those changes into the `prepare-release-prs/x` branch. You do not need
     to wait for the CI to pass in this prep-to-prep PR.

2. Create a branch named `release/<release-series>` (e.g. `release/v0.55.x`) in
   Contrib from the changelog update commit and push it to
   `open-telemetry/opentelemetry-collector-contrib`.

3. Make sure you are on `release/<release-series>`. Tag all the module groups
   (`contrib-base`) with the new release version by running the
   `make push-tags MODSET=contrib-base` command. If you set your remote using
   `https` you need to include
   `REMOTE=https://github.com/open-telemetry/opentelemetry-collector-contrib.git`
   in each command. Wait for the new tag build to pass successfully.

4. A new `v0.55.0` release should be automatically created on Github by now.
   Edit it and use the contents from the CHANGELOG.md as the release's
   description.

## 生成工件

发布过程的最后一步是为新版本的收集器创建工件，并将映像发布到 Dockerhub。该版本的
这一部分的步骤
在[opentelemetry-collector-releases](https://github.com/open-telemetry/opentelemetry-collector-releases)
repo 中完成。

1. Update the `./distribution/**/manifest.yaml` files to include the new release
   version.

2. Update the builder version in `OTELCOL_BUILDER_VERSION` to the new release in
   the `Makefile`. While this might not be strictly necessary for every release,
   this is a good practice.

3. Create a pull request with the change and ensure the build completes
   successfully. See
   [example](https://github.com/open-telemetry/opentelemetry-collector-releases/pull/71).

4. Tag with the new release version by running the `make push-tags TAG=v0.55.0`
   command. If you set your remote using `https` you need to include
   `REMOTE=https://github.com/open-telemetry/opentelemetry-collector-releases.git`
   in each command. Wait for the new tag build to pass successfully.

5. Ensure the "Release" action passes, this will

   1. push new container images to
      https://hub.docker.com/repository/docker/otel/opentelemetry-collector

   2. create a Github release for the tag and push all the build artifacts to
      the Github release. See
      [example](https://github.com/open-telemetry/opentelemetry-collector-releases/actions/runs/1346637081).

## 故障排除

1. `unknown revision internal/coreinternal/v0.55.0` -- This is typically an
   indication that there's a dependency on a new module. You can fix it by
   adding a new `replaces` entry to the `go.mod` for the affected module.
2. `commitChangesToNewBranch failed: invalid merge` -- This is a
   [known issue](https://github.com/open-telemetry/opentelemetry-go-build-tools/issues/47)
   with our release tooling. The current workaround is to clone a fresh copy of
   the repository and try again. Note that you may need to set up a `fork`
   remote pointing to your own fork for the release tooling to work properly.
3. `could not run Go Mod Tidy: go mod tidy failed` when running `multimod` --
   This is a
   [known issue](https://github.com/open-telemetry/opentelemetry-go-build-tools/issues/46)
   with our release tooling. The current workaround is to run `make gotidy`
   manually after the multimod tool fails and commit the result.
4. `Incorrect version "X" of "go.opentelemetry.io/collector/component" is included in "X"`
   in CI after `make update-otel` -- It could be because the make target was run
   too soon after updating Core and the goproxy hasn't updated yet. Try running
   `export GOPROXY=direct` and then `make update-otel`.
5. `error: failed to push some refs to 'https://github.com/open-telemetry/opentelemetry-collector-contrib.git'`
   during `make push-tags` -- If you encounter this error the `make push-tags`
   target will terminate without pushing all the tags. Using the output of the
   `make push-tags` target, save all the un-pushed the tags in `tags.txt` and
   then use this make target to complete the push:

   ```bash
   .PHONY: temp-push-tags
   temp-push-tags:
       for tag in `cat tags.txt`; do \
           echo "pushing tag $${tag}"; \
           git push ${REMOTE} $${tag}; \
       done;
   ```

## 错误修复版本

### bug 修复发布标准

Both `opentelemetry-collector` and `opentelemetry-collector-contrib` have very
short 2 week release cycles. Because of this, we put a high bar when considering
making a patch release, to avoid wasting engineering time unnecessarily.

When considering making a bugfix release on the `v0.N.x` release cycle, the bug
in question needs to fulfill the following criteria:

1. The bug has no workaround or the workaround is significantly harder to put in
   place than updating the version. Examples of simple workarounds are:
   - Reverting a feature gate.
   - Changing the configuration to an easy to find value.
2. The bug happens in common setups. To gauge this, maintainers can consider the
   following:
   - The bug is not specific to an uncommon platform
   - The bug happens with the default configuration or with a commonly used one
     (e.g. has been reported by multiple people)
3. The bug is sufficiently severe. For example (non-exhaustive list):
   - The bug makes the Collector crash reliably
   - The bug makes the Collector fail to start under an accepted configuration
   - The bug produces significant data loss
   - The bug makes the Collector negatively affect its environment (e.g.
     significantly affects its host machine)

We aim to provide a release that fixes security-related issues in at most 30
days since they are publicly announced; with the current release schedule this
means security issues will typically not warrant a bugfix release. An exception
is critical vulnerabilities (CVSSv3 score >= 9.0), which will warrant a release
within five business days.

The OpenTelemetry Collector maintainers will ultimately have the responsibility
to assess if a given bug or security issue fulfills all the necessary criteria
and may grant exceptions in a case-by-case basis.

### bug 修复发布过程

The following documents the procedure to release a bugfix

1. Create a pull request against the `release/<release-series>` (e.g.
   `release/v0.45.x`) branch to apply the fix.
2. Create a pull request to update version number against the
   `release/<release-series>` branch.
3. Once those changes have been merged, create a pull request to the `main`
   branch from the `release/<release-series>` branch.
4. Enable the **Merge pull request** setting in the repository's **Settings**
   tab.
5. Tag all the modules with the new release version by running the
   `make add-tag` command (e.g. `make add-tag TAG=v0.55.0`). Push them to
   `open-telemetry/opentelemetry-collector` with `make push-tag TAG=v0.55.0`.
   Wait for the new tag build to pass successfully.
6. **IMPORTANT**: The pull request to bring the changes from the release branch
   _MUST_ be merged using the **Merge pull request** method, and _NOT_ squashed
   using “**Squash and merge**”. This is important as it allows us to ensure the
   commit SHA from the release branch is also on the main branch. **Not
   following this step will cause much go dependency sadness.**
7. Once the branch has been merged, it will be auto-deleted. Restore the release
   branch via GitHub.
8. Once the patch is release, disable the **Merge pull request** setting.

## 发布时间表

| Date       | Version | Release manager |
| ---------- | ------- | --------------- |
| 2023-06-19 | v0.80.0 | @Aneurysm9      |
| 2023-07-03 | v0.81.0 | @jpkrohling     |
| 2023-07-17 | v0.82.0 | @mx-psi         |
| 2023-07-31 | v0.83.0 | @djaglowski     |
| 2023-08-14 | v0.84.0 | @dmitryax       |
| 2023-08-28 | v0.85.0 | @codeboten      |
| 2023-09-11 | v0.86.0 | @codeboten      |
| 2023-09-25 | v0.87.0 | @bogdandrutu    |
