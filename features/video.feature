Feature: Videos
  In order to download youtube videos
  As a user
  I need to be able to specify the videoId, the opening and ending videos as well as the video output path

  Scenario: Download a video
    Given I want to download the video that contains the id "BZP1rYjoBgI"
    When I send the download request with the output path "./" and the video itself as opening and ending
    Then the video should be in a processing state
    And I wait 60 seconds so the video can be processed and downloaded
    Then I verify if the video was really downloaded