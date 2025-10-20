# tests/integration/ai-agents/behavior-simulation.integration.test.py
import pytest
import asyncio
import time
from src.agents.human_behavior import HumanBehaviorAI
from src.agents.captcha_solver import CaptchaSolverAI
from src.agents.pattern_recognizer import PatternRecognizer

class TestBehaviorSimulationIntegration:
    
    @pytest.fixture
    def behavior_ai(self):
        return HumanBehaviorAI()
    
    @pytest.fixture
    def captcha_solver(self):
        return CaptchaSolverAI()
    
    @pytest.fixture
    def pattern_recognizer(self):
        return PatternRecognizer()
    
    @pytest.mark.integration
    @pytest.mark.asyncio
    async def test_complete_behavior_simulation_workflow(self, behavior_ai):
        """Test complete human behavior simulation from start to finish"""
        
        # Simulate a complete browsing session
        session_actions = []
        
        # 1. Initial mouse movement to search box
        print("🖱️ Simulating mouse movement to search...")
        search_trajectory = await behavior_ai.generate_mouse_trajectory(
            [100, 200], [400, 150]  # From random position to search box
        )
        session_actions.extend([
            {'type': 'mouse_move', 'data': point, 'timestamp': time.time() + i*0.1}
            for i, point in enumerate(search_trajectory)
        ])
        
        # 2. Typing search query with human-like patterns
        print("⌨️ Simulating search query typing...")
        search_query = "phone number lookup"
        keystrokes = await behavior_ai.simulate_typing_pattern(search_query, {
            'base_typing_speed': 0.12,
            'word_pause': 0.4,
            'error_rate': 0.03
        })
        session_actions.extend([
            {'type': 'typing', 'data': stroke, 'timestamp': time.time() + len(search_trajectory)*0.1 + i*0.15}
            for i, stroke in enumerate(keystrokes)
        ])
        
        # 3. Scrolling through results
        print("📜 Simulating results scrolling...")
        scroll_points = await behavior_ai.generate_scroll_behavior(1200, {
            'scroll_speed': 'medium',
            'read_time': 2.0
        })
        session_actions.extend([
            {'type': 'scroll', 'data': point, 'timestamp': time.time() + len(search_trajectory)*0.1 + len(keystrokes)*0.15 + i*0.5}
            for i, point in enumerate(scroll_points)
        ])
        
        # 4. Clicking on a result
        print("👆 Simulating result click...")
        click_trajectory = await behavior_ai.generate_mouse_trajectory(
            [400, 150], [450, 300]  # From search to result
        )
        session_actions.extend([
            {'type': 'mouse_move', 'data': point, 'timestamp': time.time() + len(search_trajectory)*0.1 + len(keystrokes)*0.15 + len(scroll_points)*0.5 + i*0.1}
            for i, point in enumerate(click_trajectory)
        ])
        
        # Add final click
        session_actions.append({
            'type': 'click',
            'data': {'button': 'left', 'x': 450, 'y': 300},
            'timestamp': time.time() + len(search_trajectory)*0.1 + len(keystrokes)*0.15 + len(scroll_points)*0.5 + len(click_trajectory)*0.1
        })
        
        # Validate complete session
        assert len(session_actions) > 30, "Session should have substantial activity"
        
        action_types = [a['type'] for a in session_actions]
        assert 'mouse_move' in action_types
        assert 'typing' in action_types
        assert 'scroll' in action_types
        assert 'click' in action_types
        
        # Verify timing is realistic
        timestamps = [a['timestamp'] for a in session_actions]
        total_duration = max(timestamps) - min(timestamps)
        assert total_duration > 10, "Session should take realistic time"
        assert total_duration < 120, "Session shouldn't be excessively long"
        
        print(f"✅ Complete behavior simulation: {len(session_actions)} actions over {total_duration:.1f}s")
        
        return session_actions
    
    @pytest.mark.integration
    @pytest.mark.asyncio
    async def test_ai_agent_coordination(self, behavior_ai, captcha_solver, pattern_recognizer):
        """Test coordination between multiple AI agents"""
        
        # Simulate a scenario requiring multiple AI agents
        scenario = {
            'detected_captcha': True,
            'required_behavior': 'stealth_browsing',
            'complexity': 'high'
        }
        
        agents_activated = []
        results = {}
        
        # 1. Pattern recognizer detects need for stealth
        if scenario['required_behavior'] == 'stealth_browsing':
            agents_activated.append('pattern_recognizer')
            stealth_pattern = await pattern_recognizer.generate_stealth_pattern()
            results['stealth_pattern'] = stealth_pattern
            
            # Apply stealth behavior
            behavior_config = stealth_pattern['behavior_config']
        
        # 2. Behavior AI executes stealth browsing
        agents_activated.append('behavior_ai')
        stealth_actions = await behavior_ai.generate_stealth_behavior(behavior_config)
        results['stealth_actions'] = stealth_actions
        
        # 3. CAPTCHA solver handles challenges if detected
        if scenario['detected_captcha']:
            agents_activated.append('captcha_solver')
            
            # Simulate CAPTCHA challenge
            test_captcha = {
                'type': 'image_text',
                'complexity': scenario['complexity'],
                'image_data': 'test_captcha_image'
            }
            
            captcha_solution = await captcha_solver.solve_captcha(test_captcha)
            results['captcha_solution'] = captcha_solution
        
        # Validate coordination
        assert len(agents_activated) >= 2, "Should activate multiple agents"
        assert 'behavior_ai' in agents_activated
        assert 'stealth_actions' in results
        
        if scenario['detected_captcha']:
            assert 'captcha_solver' in agents_activated
            assert 'captcha_solution' in results
        
        print(f"✅ AI agent coordination: {len(agents_activated)} agents activated")
        
        return results

class TestCaptchaSolvingIntegration:
    
    @pytest.mark.integration
    @pytest.mark.asyncio
    async def test_real_captcha_solving_workflow(self, captcha_solver):
        """Test complete CAPTCHA solving workflow with real challenges"""
        
        # Note: This test requires actual CAPTCHA challenges
        # For safety, we'll use simulated challenges in integration tests
        
        captcha_types = [
            {
                'type': 'image_text',
                'description': 'Distorted text CAPTCHA',
                'difficulty': 'medium'
            },
            {
                'type': 'image_selection', 
                'description': 'Select all squares with traffic lights',
                'difficulty': 'hard'
            },
            {
                'type': 'audio',
                'description': 'Audio CAPTCHA',
                'difficulty': 'medium'
            }
        ]
        
        results = {}
        
        for captcha_type in captcha_types:
            print(f"🔐 Testing {captcha_type['type']} CAPTCHA...")
            
            # Generate simulated CAPTCHA challenge
            challenge = await captcha_solver.generate_test_challenge(captcha_type)
            
            # Solve the challenge
            start_time = time.time()
            solution = await captcha_solver.solve_captcha(challenge)
            solve_time = time.time() - start_time
            
            # Validate solution
            assert 'success' in solution
            assert 'solution' in solution
            assert 'confidence' in solution
            
            if solution['success']:
                assert solve_time < 30, "CAPTCHA solving should complete within 30 seconds"
                assert solution['confidence'] > 0.5, "Should have reasonable confidence"
            
            results[captcha_type['type']] = {
                'success': solution['success'],
                'solve_time': solve_time,
                'confidence': solution['confidence']
            }
            
            print(f"   {captcha_type['type']}: {solution['success']} in {solve_time:.1f}s (conf: {solution['confidence']:.2f})")
        
        # Analyze overall performance
        success_rate = sum(1 for r in results.values() if r['success']) / len(results)
        avg_solve_time = sum(r['solve_time'] for r in results.values()) / len(results)
        
        print(f"📊 CAPTCHA solving summary: {success_rate:.1%} success rate, {avg_solve_time:.1f}s average")
        
        assert success_rate > 0.6, "Should solve majority of CAPTCHAs"
        assert avg_solve_time < 15, "Average solve time should be reasonable"
        
        return results
