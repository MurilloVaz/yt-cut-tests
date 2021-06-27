Feature: Cuts
  In order to create video cuts
  As a user
  I need to be able to specify the video and time to cut

  Scenario: Create a cut
    Given the video "BZP1rYjoBgI" is downloaded and ready to cut
    When I send a cut request with start "00:00:10" and end "00:00:20"
    Then the cut should be in a processing state
    And I wait 4 minutes so the cut request can be processed and cut
    Then I verify if the cut was really made 