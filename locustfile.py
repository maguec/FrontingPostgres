from locust import HttpUser, task, between

class QuickstartUser(HttpUser):

    @task
    def view_profiles(self):
        for profile_id in range(1, 100000):
            self.client.get(f"/profile/{profile_id}", name="/profile")

