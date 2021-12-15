#!/bin/python3
import sys
import argparse
import requests
import yaml

# HOST = "http://localhost:8080"
HOST = "https://linebot-nsun2l34rq-de.a.run.app"

def question(args):

  response =requests.post(f"{HOST}/api/chats/text", json={
    "text": "問題2",
  })
  question = response.json()
  print(question)
  response =requests.post(f"{HOST}/api/chats/text", json={
    "text": "對了",
  })
  ok = response.json()
  print(ok)
  response =requests.post(f"{HOST}/api/chats/text", json={
    "text": "再試一次",
    "delay": 1,
  })
  again = response.json()
  print(again)
  response =requests.post(f"{HOST}/api/chats/text", json={
    "text": "錯了",
    "nextChats": [{
      "id": again["id"],
    }],
  })
  error = response.json()
  print(error)
  response = requests.post(f"{HOST}/api/questions", json={
    "chat": {"id": question["id"]},
    "answer": "A",
    "options": ["A", "B", "C", "D"],
    "ok": [{"id": ok["id"]}],
    "error": [{"id": error["id"]}],
  })
  q2 = response.json()
  print(q2)

  response =requests.post(f"{HOST}/api/chats/text", json={
    "text": "問題1",
  })
  question = response.json()
  print(question)
  response =requests.post(f"{HOST}/api/chats/text", json={
    "text": "對了",
  })
  ok = response.json()
  print(ok)
  response =requests.post(f"{HOST}/api/chats/text", json={
    "text": "錯了",
  })
  error = response.json()
  print(error)
  response = requests.post(f"{HOST}/api/questions", json={
    "chat": {"id": question["id"]},
    "answer": "A",
    "options": ["A", "B", "C", "D"],
    "ok": [{"id": ok["id"]}],
    "error": [{"id": error["id"]}],
    "next" : [{"id": q2["id"]}],
  })
  if response.status_code != 200:
    print(response.text)
  print(response.json())

def loadQuestion(args):
  config = args.config[0]
  with open(config, 'r') as input:
    questionSet = [
      q for q in yaml.safe_load_all(input) if q is not None
    ]
  for questions in questionSet:
    postQuestionSet(questions)

def postQuestionSet(questions):
  questions = questions[::-1]
  next = None
  for question in questions:
    next = postQuestion(question, next)

def postQuestion(question, next):
  chat = postAction(question.get('chat', []))
  ok = [{"id": postAction(act)} for act in question.get('ok', [])]
  error = [{"id": postAction(act)} for act in question.get('error', [])]
  response = requests.post(f"{HOST}/api/questions", json={
    "chat": {"id": chat},
    "answer": question["answer"],
    "ok": ok,
    "error": error,
    "next" : [{"id": next} ] if next is not None else [],
  })
  if response.status_code != 200:
    raise RuntimeError(response.text)
  q = response.json()
  print(q)
  return q["id"]

def loadChat(args):
  config = args.config[0]
  with open(config, 'r') as input:
    actions = [
      action for action in yaml.safe_load_all(input)
      if action is not None
    ]
  for action in actions:
    postAction(action)

def postAction(action) -> int:
  chats = action[::-1]
  nxt: int = None
  for chat in chats:
    nxt = postChat(chat, nxt)
  return nxt

def postChat(chat, nxt: int = None) -> int:
  msgType = None
  if "text" in chat:
    msgType = "text"
  elif "sticker" in chat:
    msgType = "sticker"
  elif "image" in chat:
    msgType = "image"
  if nxt:
    chat["nextChats"] = [{"id": nxt}]
  response =requests.post(f"{HOST}/api/chats/{msgType}", json=chat)
  if response.status_code != 200:
    raise RuntimeError(response.text)
  result = response.json()
  print(result)
  return result["id"]

def main():
  parser = argparse.ArgumentParser(description="linebot util")
  subparsers = parser.add_subparsers()

  test_cmd = subparsers.add_parser("question")
  test_cmd.set_defaults(func=question)

  load_chat_cmd = subparsers.add_parser("loadChat")
  load_chat_cmd.set_defaults(func=loadChat)
  load_chat_cmd.add_argument("config", nargs=1, type=str)

  load_question_cmd = subparsers.add_parser("loadQuestion")
  load_question_cmd.set_defaults(func=loadQuestion)
  load_question_cmd.add_argument("config", nargs=1, type=str)

  args = parser.parse_args(sys.argv[1:])
  args.func(args)


if __name__ == "__main__":
  main()
