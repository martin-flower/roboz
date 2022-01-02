Feature: roboz entry endpoint
  In order to clean the office
  As a robot controller
  I need to be able to send instructions to a robot
  I need to verify that instructions have been completed  

  Scenario: invalid commands
    When the controller sends invalid commands to the robot
    Then an error response is received

  Scenario: clean the office with one instruction
    When the controller sends 1 instructions to the robot
    Then a response is received which matches the instructions

  Scenario: clean the office with 10 instructions
    When the controller sends 10 instructions to the robot
    When the controller asks for a list of executions
    Then the number of commands is at least 10
