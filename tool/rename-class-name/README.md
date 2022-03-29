# Update class name for testing repot
* It seems unfortunately the JUnit testing report converted from golang [testing result](https://github.com/jstemmer/go-junit-report) just simply uses the directory name as package which will be considered as "." package in Java context

# How to use it
* go install 
* rename-class-name -file <your-junit-report.xml> -server <your-desired-package-name>, 
    * e.g, rename-class-name -file deduction-report.xml -server paidads.deduction
