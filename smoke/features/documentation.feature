Feature: roboz documentation
  In order to understand roboz
  As a robot controller
  I need to be able to read the documentation

  Scenario: read the documentation
    When the controller asks to read the documentation
    Then the documentation page is displayed
    