# micro-ai-agents/python/src/agents/human_behavior.py
import numpy as np
from typing import List, Dict
import torch
import torch.nn as nn

class HumanBehaviorAI:
    def __init__(self):
        self.mouse_model = self.load_mouse_model()
        self.typing_model = self.load_typing_model()
        self.behavior_predictor = self.load_behavior_predictor()
    
    def generate_mouse_trajectory(self, start_pos: List[float], end_pos: List[float]) -> List[Dict]:
        """Generate human-like mouse movement"""
        points = []
        current = np.array(start_pos)
        target = np.array(end_pos)
        
        # Bezier curve with human-like variations
        for t in np.linspace(0, 1, 20):
            # Control points with randomness
            control1 = current + (target - current) * 0.3 + np.random.normal(0, 5, 2)
            control2 = current + (target - current) * 0.7 + np.random.normal(0, 5, 2)
            
            point = self.bezier_point(current, control1, control2, target, t)
            
            # Add human jitter
            jitter = np.random.normal(0, 1.5, 2)
            point += jitter
            
            points.append({
                'x': float(point[0]),
                'y': float(point[1]),
                'time_offset': t * np.random.uniform(0.8, 1.2)
            })
        
        return points
    
    def simulate_typing_pattern(self, text: str, persona: Dict) -> List[Dict]:
        """Generate realistic typing rhythm"""
        keystrokes = []
        words = text.split()
        
        for word in words:
            # Type word with character delays
            for char in word:
                delay = persona['base_typing_speed'] * np.random.uniform(0.8, 1.2)
                keystrokes.append({
                    'char': char,
                    'delay': delay,
                    'timestamp': sum(k['delay'] for k in keystrokes)
                })
            
            # Word pause
            pause = persona['word_pause'] * np.random.uniform(0.7, 1.3)
            keystrokes.append({
                'char': 'PAUSE',
                'delay': pause,
                'timestamp': sum(k['delay'] for k in keystrokes)
            })
        
        return keystrokes

class CaptchaSolverAI:
    def __init__(self):
        self.text_model = self.load_text_recognition()
        self.image_model = self.load_image_analysis()
    
    def solve_text_captcha(self, image_data) -> str:
        """Solve text-based CAPTCHAs using our AI"""
        # Preprocess image
        processed = self.preprocess_captcha(image_data)
        
        # Use ensemble of models
        predictions = []
        for model in [self.text_model, self.backup_model]:
            try:
                prediction = model.predict(processed)
                predictions.append(prediction)
            except Exception as e:
                continue
        
        return self.consensus_prediction(predictions)
    
    def solve_image_captcha(self, images, question) -> int:
        """Solve image selection CAPTCHAs"""
        # Analyze each image
        scores = []
        for img in images:
            relevance = self.analyze_image_relevance(img, question)
            scores.append(relevance)
        
        return np.argmax(scores)
