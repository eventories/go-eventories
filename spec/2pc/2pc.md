# 2-phase-commit
go-eventories is works based on leader/follower model for database syncing. <br>

# prepare
<img width="485" alt="normal" src="https://user-images.githubusercontent.com/72970043/210774211-59cda4cc-5480-4829-b74e-1c526f3e7513.png">

<img width="505" alt="leader_prepare_failed" src="https://user-images.githubusercontent.com/72970043/210774300-34144ff0-54ec-4e2c-85b7-f1c3b6e3615d.png">

<img width="495" alt="leader_send_prepare_failed" src="https://user-images.githubusercontent.com/72970043/210774318-9328a909-e577-4614-ae91-1701969f2c37.png">

<img width="492" alt="follower_prepare_failed" src="https://user-images.githubusercontent.com/72970043/210774360-de5aaa58-1420-4e75-b770-060d64cb6f24.png">

<img width="499" alt="follower_send_ack_failed" src="https://user-images.githubusercontent.com/72970043/210774370-96e3e954-5f2d-4c43-84bf-9603a170e865.png">

<img width="495" alt="leader_send_ack_failed" src="https://user-images.githubusercontent.com/72970043/210774391-c8d2947a-d41c-4f22-a17f-baf7f778b03e.png">


# commit

<img width="481" alt="normal" src="https://user-images.githubusercontent.com/72970043/210775405-384bf1ff-b567-40d2-ad3e-e1a49b0b66eb.png">

<img width="494" alt="leader_commit_failed" src="https://user-images.githubusercontent.com/72970043/210775465-2299696e-4751-4c67-b853-b0b8201cc214.png">

<img width="500" alt="leader_send_ack_failed" src="https://user-images.githubusercontent.com/72970043/210775484-4fd48c9d-1e7d-4ff9-a322-8e1cfeb0f39a.png">

<img width="489" alt="follower_commit_failed" src="https://user-images.githubusercontent.com/72970043/210775507-01d755d3-3114-494f-a058-470a7f8fa8f0.png">

<img width="500" alt="follower_send_ack_failed" src="https://user-images.githubusercontent.com/72970043/210775519-672a253a-efcf-4f1b-9e4f-e44311f6c8e5.png">

<img width="512" alt="leader_send_ack2_failed" src="https://user-images.githubusercontent.com/72970043/210775545-fdba19a3-fdf4-4af3-93bd-cec0d995db84.png">
