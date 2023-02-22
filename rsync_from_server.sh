echo 'rsync --checksum -av --exclude-from=rsync_exclude.txt foo@foo.sakura.ne.jp:~/"*bar*" ./'
echo
rsync -n --checksum -av --exclude-from=rsync_exclude.txt foo@foo.sakura.ne.jp:~/"*bar*" ./
