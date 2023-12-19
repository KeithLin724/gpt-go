import ast

import logging

from fastapi import FastAPI, Response, Request
from fastapi.responses import StreamingResponse
from typing import List, Union, Any, Dict, AnyStr

# from ._tokenizer        import tokenize
# from .. import BaseProvider

import time
import json
import random
import string
import uvicorn
import nest_asyncio

import g4f

import g4f
import asyncio
from enum import Enum

# https://myapollo.com.tw/blog/begin-to-asyncio/


class TextGenerator:
    """
    The `TextGenerator` class is a Python class that generates text using the GPT-4 model and multiple
    providers asynchronously."""

    class MessageState(Enum):
        ok = "OK"
        err = "ERR"

    G4F_VERSION = g4f.version

    def __init__(self) -> None:
        self._provide = [
            g4f.Provider.Bing,
            g4f.Provider.ChatBase,
            g4f.Provider.GptGo,
            g4f.Provider.You,
            g4f.Provider.Raycast,
            # g4f.Provider.Yqcloud,
        ]

        self._cut_prompt_provider = [
            g4f.Provider.You,
            # g4f.Provider.GptGo,
            g4f.Provider.GeekGpt,
        ]

    async def run_provider(self, provider: g4f.Provider.BaseProvider, prompt: str):
        """
        The function `run_provider` takes a provider and a prompt as input, and uses the provider to
        generate a response to the prompt using the GPT-4 model.

        :param provider: The `provider` parameter is an instance of a class that inherits from
        `g4f.Provider.BaseProvider`. It is used to specify the provider for the chat completion model
        :type provider: g4f.Provider.BaseProvider
        :param prompt: The `prompt` parameter is a string that represents the user's input or message to the
        chatbot. It is the content that the user wants to send to the chatbot for processing
        :type prompt: str
        :return: The function `run_provider` returns a tuple containing three elements:
        """
        try:
            response = await g4f.ChatCompletion.create_async(
                model=g4f.models.gpt_4,
                messages=[
                    {
                        "role": "user",
                        "content": prompt,
                    }
                ],
                provider=provider,
                timeout=5,
            )
            # print(f"{provider.__name__}:", response)
            return "OK", provider.__name__, response
        except Exception as e:
            return "ERR", provider.__name__, str(e)

    async def get_waiting(self, task_list: list):
        """
        The `get_waiting` function waits for the first task in a list to complete and returns its result
        along with the remaining pending tasks.

        :param task_list: The `task_list` parameter is a list of tasks that you want to wait for. These
        tasks are typically created using the `asyncio.create_task()` function or by wrapping a coroutine
        function with `asyncio.ensure_future()`. Each task represents a concurrent operation that you want
        to wait for
        :type task_list: list
        :return: The function `get_waiting` returns two values: `result` and `pending_tasks`.
        """
        try:
            # Wait for the first task to complete
            completed_task, pending_tasks = await asyncio.wait(
                task_list,
                return_when=asyncio.FIRST_COMPLETED,
            )

            # Check the result of the first completed task
            result = await completed_task.pop()

            return result, pending_tasks
        except asyncio.CancelledError as e:
            print(e)
            pass

    async def get_generate(self, prompt: str):
        """
        The function `get_generate` takes a prompt as input and runs multiple providers asynchronously,
        returning the first successful result or the first error encountered.

        :param prompt: The `prompt` parameter is a string that represents the input prompt for
        generating text. It is used as an input for each provider in the `_provide` list
        :type prompt: str
        :return: a tuple containing the provider name and the generated message. If all tasks encounter
        an error, the function returns the first error encountered. If there are no pending tasks, an
        empty string is returned.
        """
        pending_tasks = [
            asyncio.create_task(
                self.run_provider(
                    provider=provider,
                    prompt=prompt,
                )
            )
            for provider in self._provide
        ]

        first_result = None
        # result, pending_tasks = await self.get_waiting(pending_tasks)
        for _ in range(len(pending_tasks)):
            result, pending_tasks = await self.get_waiting(pending_tasks)
            state, provider_name, msg = result

            if pending_tasks is None:
                return ""

            # is ok
            if state == "OK" and msg != "":
                # if some task is not done, cancel it
                for task in pending_tasks:
                    if not task.done():
                        task.cancel()

                # [task.cancel() for task in pending_tasks]

                return provider_name, msg

            # handel error
            elif state == "ERR":
                first_result = result

        return first_result[-1]

    async def generate(self, prompt):
        """
        The function "generate" is an asynchronous function that takes a prompt as input and returns the
        result of calling the "get_generate" function with the prompt as an argument.

        :param prompt: The prompt is the input text or sentence that you want to use as a starting point
        for generating the output
        :return: The result of the `get_generate` method is being returned.
        """

        provider, msg = await self.get_generate(prompt=prompt)
        return provider, msg


LOGGER = logging.getLogger("uvicorn")
text_generator = TextGenerator()


class Api:
    def __init__(
        self,
        engine: g4f,
        debug: bool = True,
        sentry: bool = False,
        # list_ignored_providers: List[Union[str, None]] = None,
    ) -> None:
        self.engine = engine
        self.debug = debug
        self.sentry = sentry
        # self.list_ignored_providers = list_ignored_providers
        # self.list_ignored_providers = [
        #     g4f.Provider.Raycast,
        #     g4f.Provider.Phind,
        #     g4f.Provider.Bing,
        #     # g4f.Provider.Liaobots,
        # ]

        self.app = FastAPI()
        nest_asyncio.apply()

        JSONObject = Dict[AnyStr, Any]
        JSONArray = List[Any]
        JSONStructure = Union[JSONArray, JSONObject]

        @self.app.get("/")
        async def read_root():
            return Response(
                content=json.dumps({"info": "g4f API"}, indent=4),
                media_type="application/json",
            )

        @self.app.get("/v1")
        async def read_root_v1():
            return Response(
                content=json.dumps(
                    {"info": "Go to /v1/chat/completions or /v1/models."}, indent=4
                ),
                media_type="application/json",
            )

        @self.app.get("/v1/models")
        async def models():
            model_list = []
            for model in g4f.Model.__all__():
                model_info = g4f.ModelUtils.convert[model]
                model_list.append(
                    {
                        "id": model,
                        "object": "model",
                        "created": 0,
                        "owned_by": model_info.base_provider,
                    }
                )
            return Response(
                content=json.dumps({"object": "list", "data": model_list}, indent=4),
                media_type="application/json",
            )

        @self.app.get("/v1/models/{model_name}")
        async def model_info(model_name: str):
            try:
                model_info = g4f.ModelUtils.convert[model_name]

                return Response(
                    content=json.dumps(
                        {
                            "id": model_name,
                            "object": "model",
                            "created": 0,
                            "owned_by": model_info.base_provider,
                        },
                        indent=4,
                    ),
                    media_type="application/json",
                )
            except:
                return Response(
                    content=json.dumps(
                        {"error": "The model does not exist."}, indent=4
                    ),
                    media_type="application/json",
                )

        @self.app.post("/v1/chat/completions")
        async def chat_completions(request: Request, item: JSONStructure = None):
            item_data = {
                "model": "gpt-3.5-turbo",
                "stream": False,
            }

            # item contains byte keys, and dict.get suppresses error
            item_data.update(
                {
                    key.decode("utf-8") if isinstance(key, bytes) else key: str(value)
                    for key, value in (item or {}).items()
                }
            )

            # messages is str, need dict
            if isinstance(item_data.get("messages"), str):
                item_data["messages"] = ast.literal_eval(item_data.get("messages"))

            model = item_data.get("model")
            stream = True if item_data.get("stream") == "True" else False
            messages = item_data.get("messages")
            conversation = (
                item_data.get("conversation")
                if item_data.get("conversation") != None
                else None
            )

            provider = None

            try:
                provider, response = await text_generator.generate(messages)
                LOGGER.info(f"{provider}: {response}")
            except Exception as e:
                LOGGER.exception(e)
                return {
                    "error": "An error occurred while generating the response.",
                }

            completion_id = "".join(
                random.choices(string.ascii_letters + string.digits, k=28)
            )
            completion_timestamp = int(time.time())

            if not stream:
                # prompt_tokens, _ = tokenize(''.join([message['content'] for message in messages]))
                # completion_tokens, _ = tokenize(response)

                return {
                    "id": f"chatcmpl-{completion_id}",
                    "object": "chat.completion",
                    "created": completion_timestamp,
                    "model": model,
                    "choices": [
                        {
                            "index": 0,
                            "message": {
                                "role": "assistant",
                                "content": response,
                            },
                            "finish_reason": "stop",
                        }
                    ],
                    "usage": {
                        "prompt_tokens": 0,  # prompt_tokens,
                        "completion_tokens": 0,  # completion_tokens,
                        "total_tokens": 0,  # prompt_tokens + completion_tokens,
                    },
                }

            def streaming():
                try:
                    for chunk in response:
                        completion_data = {
                            "id": f"chatcmpl-{completion_id}",
                            "object": "chat.completion.chunk",
                            "created": completion_timestamp,
                            "model": model,
                            "choices": [
                                {
                                    "index": 0,
                                    "delta": {
                                        "role": "assistant",
                                        "content": chunk,
                                    },
                                    "finish_reason": None,
                                }
                            ],
                        }

                        content = json.dumps(completion_data, separators=(",", ":"))
                        yield f"data: {content}\n\n"
                        time.sleep(0.03)

                    end_completion_data = {
                        "id": f"chatcmpl-{completion_id}",
                        "object": "chat.completion.chunk",
                        "created": completion_timestamp,
                        "model": model,
                        "choices": [
                            {
                                "index": 0,
                                "delta": {},
                                "finish_reason": "stop",
                            }
                        ],
                    }

                    content = json.dumps(end_completion_data, separators=(",", ":"))
                    yield f"data: {content}\n\n"

                except GeneratorExit:
                    pass

            return StreamingResponse(streaming(), media_type="text/event-stream")

        @self.app.post("/v1/completions")
        async def completions():
            return Response(
                content=json.dumps({"info": "Not working yet."}, indent=4),
                media_type="application/json",
            )

    def run(self, ip):
        split_ip = ip.split(":")
        uvicorn.run(
            app=self.app, host=split_ip[0], port=int(split_ip[1]), use_colors=False
        )


Api(g4f).run("0.0.0.0:1337")
