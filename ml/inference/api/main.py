# ml/inference/api/main.py
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

class PredictionRequest(BaseModel):
    phone_number: str
    features: dict

@app.post("/predict")
async def predict_behavior(request: PredictionRequest):
    model = load_model("models/behavior_v1.pkl")
    prediction = model.predict([request.features])
    return {"risk_score": prediction[0]}
