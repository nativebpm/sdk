import unittest
from nativebpm.builder import Workflow

class TestBuilderBranching(unittest.TestCase):
    def test_when_then_otherwise(self):
        w = Workflow("approval_process", "Approval Process")
        (w.start("start")
         .service("check_risk", "Check Risk", "risk_eval")
         .when("risk == 'low'")
         .then(lambda b: b.service("auto_approve", "Auto Approve", "fast_track"))
         .when("risk == 'medium'")
         .then(lambda b: b.user("manager_review", "Manager Review"))
         .otherwise(lambda b: b.service("auto_reject", "Auto Reject", "reject_topic"))
         .end("end", "Done"))

        nodes = {n["id"]: n for n in w._nodes}
        self.assertIn("start", nodes)
        self.assertIn("check_risk", nodes)
        self.assertIn("auto_approve", nodes)
        self.assertIn("manager_review", nodes)
        self.assertIn("auto_reject", nodes)
        self.assertIn("end", nodes)

if __name__ == "__main__":
    unittest.main()
