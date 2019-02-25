## Sample thundr application

The thundr-sample project contains sample implementations of the thundr
framework for different deployment environments.

You can read more about thundr [here](http://3wks.github.com/thundr).

You can see this sample in action [here](http://gae.thundr-sample.appspot.com).

### Getting it
You can obtain the source by cloning the git repository from github:

	git clone https://github.com/3wks/thundr-sample.git
	
Each branch is for a different employment environment. You can list the examples
available by using, 
	
	git branch
	
Switch to a particular example by checking out the branch
		
	git checkout <branch-name>
	e.g. git checkout gae
	
You can also use the samples as scaffolding - just checkout the branch you want,
remove the git repository and delete any sample code you're not interested in.

	e.g.
	git checkout gae
	rm -rf .git
	
### thundr-gae

This branch contains a working implementation of a thundr application
running using Google's AppEngine (GAE). It relies on the thundr-gae module.

Read more about thundr-gae [here](http://3wks.github.io/thundr/thundr-gae/index.html)

To run, use the appengine maven plugin:

	mvn appengine:devserver
	
To deploy to appengine, update the <application> element in /src/main/webapp/WEB-INF/appengine-web.xml to an app-id you own,
then run:

    mvn appengine:update

--------------    
thundr sample application - Copyright (C) 2013 3wks
