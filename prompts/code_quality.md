# Expectations for Code Quality

All the following are expected in a response, most of these are not
hard requirements but are expected to be present in a good solution.

- Readability
- Good names
- Good error handling
- DRY principles
- SOLID principles
- Separation of concerns  
- Testability
- Maintainability
- Scalability
- Do not use comments unless strictly necessary

Write me some top tier code. The best code. 
Write me some code that will last a thousand years.

If you have to refactor some surrounding logic
in order to accomplish a better solution, feel free to
do this.

Follow existing patterns in the codebase.

# Unit test expectations

When writing unit test suites, be sure to 
create helper functions which build the test 
data for you and reference these in every test.

Each individual test should be very simple and
readable, ideally the names are so well defined
it could be read by a non-technical person.

For writing a test suite for a package, you should
define a base test suite which handles all of the 
setup and teardown for the tests, including database
connections etc. 

Other test suites should inherit from this
    
Unit tests should not have comments. Do not comment
unit tests, the test name should be descriptive enough
