# Xena - the Amazon Warrior Princess

Just a command line tool to assist automation of AEM and AWS in general

So far it only displays
* EC2 instances with a tag of Role and Environment
* latest ebs snapshot matching a name

```
./xena snapshots --name aem-author --latest
snap-8975ee86,aem-author-test-snap-20160314-1657,2016-03-14 05:57:51 +0000 UTC
```

```
./xena snapshots --name aem-author --latest --summary
snap-8975ee86
```

show all snapshots containing name

```
./myaws snapshots --name aem-author 
snap-8975ee86,aem-author-test-snap-20160314-1657,2016-03-14 05:57:51 +0000 UTC
```

