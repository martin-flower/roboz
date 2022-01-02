Feature: roboz health
  In order to know whether roboz can be used
  As a robot controller
  I need to be know if roboz is healthy

  Scenario: roboz is healthy
    When the controller asks roboz whether it is healthy
    Then roboz replies ok
    