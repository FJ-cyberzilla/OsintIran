# tests/integration/ai_agents/behavior_simulation.test.py
import pytest
import asyncio
from unittest.mock import Mock, patch
from src.agents.human_behavior import HumanBehaviorAI
from src.agents.captcha_solver import CaptchaSolverAI

class TestHumanBehaviorAI:
    
    @pytest.fixture
    def behavior_ai(self):
        return HumanBehaviorAI()
    
    @pytest.fixture
    def captcha_solver(self):
        return CaptchaSolverAI()
    
    @pytest.mark.asyncio
    async def test_mouse_trajectory_generation(self, behavior_ai):
        """Test that mouse trajectories are generated realistically"""
        start_pos = [100, 100]
        end_pos = [500, 500]
        
        trajectory = await behavior_ai.generate_mouse_trajectory(start_pos, end_pos)
        
        assert len(trajectory) > 10, "Trajectory should have multiple points"
        assert trajectory[0]['x'] == start_pos[0]
        assert trajectory[0]['y'] == start_pos[1]
        assert trajectory[-1]['x'] == end_pos[0]
        assert trajectory[-1]['y'] == end_pos[1]
        
        # Check that points follow a natural curve
        for i in range(1, len(trajectory)):
            prev_point = trajectory[i-1]
            curr_point = trajectory[i]
            assert curr_point['time_offset'] > prev_point['time_offset']
    
    @pytest.mark.asyncio
    async def test_typing_pattern_simulation(self, behavior_ai):
        """Test realistic typing pattern generation"""
        test_text = "Hello World"
        persona = {
            'base_typing_speed': 0.1,
            'word_pause': 0.5,
            'error_rate': 0.05
        }
        
        keystrokes = await behavior_ai.simulate_typing_pattern(test_text, persona)
        
        assert len(keystrokes) >= len(test_text), "Should have at least one keystroke per character"
        
        total_time = sum(k['delay'] for k in keystrokes)
        assert total_time > 0, "Should have positive total typing time"
        
        # Check for natural pauses between words
        word_pauses = [k for k in keystrokes if k['char'] == 'PAUSE']
        assert len(word_pauses) >= 1, "Should have pauses between words"
    
    @pytest.mark.asyncio
    async def test_scroll_behavior_generation(self, behavior_ai):
        """Test realistic scroll behavior simulation"""
        page_height = 2000
        persona = {
            'scroll_speed': 'medium',
            'read_time': 2.0
        }
        
        scroll_points = await behavior_ai.generate_scroll_behavior(page_height, persona)
        
        assert len(scroll_points) > 0, "Should generate scroll points"
        
        total_scroll = sum(point['amount'] for point in scroll_points if point['amount'] > 0)
        assert total_scroll >= page_height, "Should scroll through entire page"
        
        # Check for natural pauses
        pauses = [point for point in scroll_points if point['amount'] == 0]
        assert len(pauses) > 0, "Should include reading pauses"

class TestCaptchaSolverAI:
    
    @pytest.mark.asyncio
    async def test_text_captcha_solving(self, captcha_solver):
        """Test text-based CAPTCHA solving"""
        # Mock CAPTCHA image
        mock_image = Mock()
        mock_image.size = (200, 80)
        
        with patch.object(captcha_solver.text_model, 'predict') as mock_predict:
            mock_predict.return_value = "ABCD123"
            
            solution = await captcha_solver.solve_text_captcha(mock_image)
            
            assert solution == "ABCD123"
            mock_predict.assert_called_once()
    
    @pytest.mark.asyncio
    async def test_image_captcha_solving(self, captcha_solver):
        """Test image selection CAPTCHA solving"""
        test_images = [Mock(), Mock(), Mock(), Mock()]
        question = "Select all images with cars"
        
        with patch.object(captcha_solver.image_model, 'analyze_image_relevance') as mock_analyze:
            mock_analyze.side_effect = [0.1, 0.9, 0.2, 0.8]
            
            selected_indices = await captcha_solver.solve_image_captcha(test_images, question)
            
            assert selected_indices == [1, 3]  # Indices with highest relevance
            assert mock_analyze.call_count == len(test_images)
    
    @pytest.mark.asyncio
    async def test_captcha_fallback_mechanisms(self, captcha_solver):
        """Test CAPTCHA solving fallback when primary method fails"""
        mock_image = Mock()
        
        # Make primary model fail
        with patch.object(captcha_solver.text_model, 'predict', side_effect=Exception("Model error")):
            with patch.object(captcha_solver.backup_model, 'predict') as mock_backup:
                mock_backup.return_value = "FALLBACK"
                
                solution = await captcha_solver.solve_text_captcha(mock_image)
                
                assert solution == "FALLBACK"
                mock_backup.assert_called_once()

@pytest.mark.integration
class TestAIAgentIntegration:
    """Integration tests for AI agents working together"""
    
    @pytest.mark.asyncio
    async def test_complete_behavior_simulation(self):
        """Test complete human behavior simulation workflow"""
        behavior_ai = HumanBehaviorAI()
        
        # Simulate complete browsing session
        actions = []
        
        # Generate mouse movement
        mouse_trajectory = await behavior_ai.generate_mouse_trajectory([100, 100], [300, 300])
        actions.extend([{'type': 'mouse_move', 'data': point} for point in mouse_trajectory])
        
        # Generate typing
        keystrokes = await behavior_ai.simulate_typing_pattern("search query", {
            'base_typing_speed': 0.08,
            'word_pause': 0.3,
            'error_rate': 0.02
        })
        actions.extend([{'type': 'typing', 'data': stroke} for stroke in keystrokes])
        
        # Generate scrolling
        scroll_points = await behavior_ai.generate_scroll_behavior(1500, {
            'scroll_speed': 'fast',
            'read_time': 1.5
        })
        actions.extend([{'type': 'scroll', 'data': point} for point in scroll_points])
        
        assert len(actions) > 20, "Should generate comprehensive behavior sequence"
        
        # Verify action sequence is realistic
        action_types = [a['type'] for a in actions]
        assert 'mouse_move' in action_types
        assert 'typing' in action_types  
        assert 'scroll' in action_types
        
        # Verify timing is natural
        total_duration = sum(
            a['data'].get('time_offset', 0) + a['data'].get('delay', 0) 
            for a in actions if 'data' in a
        )
        assert total_duration > 5, "Session should have realistic duration"
