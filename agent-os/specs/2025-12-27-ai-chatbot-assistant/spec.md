# AI Chatbot Assistant - Specification

**Version:** 1.0.0  
**Date:** December 27, 2025  
**Priority:** P1 - High  
**Estimated Duration:** 4-6 weeks  

---

## Executive Summary

An intelligent AI-powered chatbot that provides 24/7 support, guides users through registration processes, answers compliance questions, and offers personalized recommendations. This feature significantly improves user experience and reduces support costs by 60%.

---

## Business Case

### Market Need
- Users need instant answers to compliance questions
- After-hours support is expensive
- Manual support doesn't scale
- Complex registration processes need guidance
- Competitors lack intelligent assistance

### Business Impact
- **Support Cost Reduction**: 60% reduction in support tickets
- **User Satisfaction**: 45% improvement in CSAT
- **Conversion Rate**: 30% increase with guided workflows
- **24/7 Availability**: No additional staff costs
- **Multilingual Support**: Serve all 11 SA official languages
- **Competitive Advantage**: First intelligent compliance assistant in SADC

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    FRONTEND                              │
│  (Chat Widget + Embedded in all pages)                   │
└────────────────────────┬─────────────────────────────────┘
                         │ WebSocket + REST
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  API GATEWAY                             │
│  (Authentication + Rate Limiting)                        │
└────────────────────────┬─────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│            CHATBOT SERVICE (Go/Node.js)                  │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │  WebSocket   │───▶│   Intent     │                  │
│  │   Server     │    │ Recognition  │                  │
│  └──────────────┘    └──────┬───────┘                  │
│                              │                           │
│                              ▼                           │
│                      ┌──────────────┐                   │
│                      │  LLM Engine  │                   │
│                      │  (OpenAI)    │                   │
│                      └──────┬───────┘                   │
│                              │                           │
│                ┌─────────────┼─────────────┐            │
│                │             │             │            │
│                ▼             ▼             ▼            │
│         ┌──────────┐  ┌──────────┐  ┌──────────┐      │
│         │Knowledge │  │ Actions  │  │Analytics │      │
│         │   Base   │  │  Engine  │  │  Engine  │      │
│         └──────────┘  └──────────┘  └──────────┘      │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  DATA LAYER                              │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │  PostgreSQL  │    │    Redis     │                  │
│  │(Conversations)│    │   (Cache)    │                  │
│  └──────────────┘    └──────────────┘                  │
└─────────────────────────────────────────────────────────┘
```

---

## Core Features

### 1. Conversational AI

**Natural Language Understanding:**
- Intent recognition (greeting, question, complaint, request)
- Entity extraction (company name, registration number, dates)
- Context management (multi-turn conversations)
- Sentiment analysis
- Language detection (11 SA languages)

**LLM Integration:**
- **Primary**: OpenAI GPT-4 Turbo
- **Fallback**: Anthropic Claude 3
- **Local Option**: Llama 3 (for sensitive data)

**Capabilities:**
- Answer compliance questions
- Explain registration processes
- Guide through forms step-by-step
- Provide document checklists
- Calculate fees
- Check status of applications
- Book appointments
- Transfer to human agents

---

### 2. Knowledge Base

**Content Categories:**
- Company registration guides
- Tax compliance information
- CIPC regulations
- SARS requirements
- B-BBEE certification info
- CIDB registration process
- CSD vendor registration
- Import/export licensing
- Trademark registration
- FAQs (500+ questions)

**Knowledge Sources:**
- Internal documentation
- Government websites
- Legal resources
- Historical support tickets
- User manuals
- Case studies

**RAG (Retrieval Augmented Generation):**
```
User Query
  ↓
Vector Embedding
  ↓
Semantic Search (Pinecone/Weaviate)
  ↓
Retrieve Top 5 Documents
  ↓
Pass to LLM with Context
  ↓
Generate Response
```

---

### 3. Guided Workflows

**Interactive Guides:**
- Company registration wizard
- Tax registration walkthrough
- Document upload assistant
- Payment processing helper
- Form completion guide

**Smart Suggestions:**
- Next step recommendations
- Missing document alerts
- Optimal service suggestions
- Cost optimization tips
- Timeline estimates

**Example Conversation:**
```
User: I want to register a company

Bot: Great! I'll help you register your company. Let me ask a few questions:

1. What type of company? 
   - Private Company (Pty Ltd)
   - Personal Liability Company (Inc)
   - Non-Profit Company (NPC)

User: Private Company

Bot: Perfect! For a Pty Ltd, you'll need:
✓ Company name (3 options)
✓ Registered address
✓ Directors (minimum 1)
✓ Shareholders
✓ ID documents

Estimated time: 7-14 days
Cost: R175 + optional services

Shall we start with choosing a company name?

User: Yes

Bot: Please provide 3 company name options in order of preference...
```

---

### 4. Proactive Assistance

**Contextual Help:**
- Detect user struggles (e.g., spending >2 min on a field)
- Offer help automatically
- Suggest similar successful cases
- Provide inline explanations

**Smart Notifications:**
- Remind about pending documents
- Alert about expiring licenses
- Notify about new regulations
- Suggest renewal dates

**Personalized Recommendations:**
- Based on user's business type
- Based on previous registrations
- Based on industry trends
- Based on compliance needs

---

### 5. Multilingual Support

**Supported Languages:**
- English
- Afrikaans
- isiZulu
- isiXhosa
- Sesotho
- Setswana
- Sepedi
- Xitsonga
- siSwati
- Tshivenda
- isiNdebele

**Implementation:**
- Automatic language detection
- Real-time translation
- Culturally appropriate responses
- Local terminology

---

### 6. Human Handoff

**Escalation Triggers:**
- User requests human agent
- Bot confidence < 70%
- Sensitive issues (complaints, refunds)
- Complex cases
- After 3 failed attempts

**Handoff Process:**
```
Bot: I'll connect you with a human agent who can better assist you.

[Transfer conversation history]

Agent: Hi, I'm Sarah. I can see you're asking about [topic]. 
       Let me help you with that...
```

**Agent Dashboard:**
- Active chats
- Conversation history
- User profile
- Suggested responses (AI-assisted)
- Quick replies
- Knowledge base access

---

## Technical Specifications

### Technology Stack

**Backend:**
- **Language**: Node.js (TypeScript) or Go
- **Framework**: Express/Fastify or Gin
- **WebSocket**: Socket.io or Gorilla WebSocket
- **LLM**: OpenAI API, LangChain
- **Vector DB**: Pinecone or Weaviate
- **Cache**: Redis
- **Database**: PostgreSQL

**Frontend:**
- **Framework**: React/Svelte
- **UI Library**: Custom chat widget
- **State Management**: Zustand/Svelte stores
- **WebSocket**: Socket.io client
- **Styling**: TailwindCSS

---

### System Architecture (Node.js Implementation)

**Project Structure:**
```
apps/chatbot/
├── src/
│   ├── server.ts
│   ├── websocket/
│   │   ├── handler.ts
│   │   └── events.ts
│   ├── llm/
│   │   ├── openai.ts
│   │   ├── anthropic.ts
│   │   ├── prompt-templates.ts
│   │   └── chain.ts
│   ├── knowledge/
│   │   ├── embeddings.ts
│   │   ├── retrieval.ts
│   │   └── indexing.ts
│   ├── intents/
│   │   ├── classifier.ts
│   │   ├── entities.ts
│   │   └── actions.ts
│   ├── actions/
│   │   ├── check-status.ts
│   │   ├── calculate-fee.ts
│   │   ├── book-appointment.ts
│   │   └── transfer-agent.ts
│   ├── conversations/
│   │   ├── manager.ts
│   │   ├── context.ts
│   │   └── history.ts
│   ├── analytics/
│   │   ├── tracker.ts
│   │   └── metrics.ts
│   └── integrations/
│       ├── comply360-api.ts
│       ├── odoo.ts
│       └── payment-gateway.ts
├── Dockerfile
├── package.json
└── tsconfig.json
```

---

### LLM Integration (LangChain)

```typescript
import { ChatOpenAI } from "@langchain/openai";
import { PromptTemplate } from "@langchain/core/prompts";
import { StringOutputParser } from "@langchain/core/output_parsers";
import { RunnableSequence } from "@langchain/core/runnables";

// Initialize LLM
const llm = new ChatOpenAI({
  modelName: "gpt-4-turbo-preview",
  temperature: 0.7,
  openAIApiKey: process.env.OPENAI_API_KEY,
});

// System prompt
const systemPrompt = `You are Comply360 Assistant, an expert in South African corporate compliance and company registration.

Your role:
- Help users with company registration and compliance
- Provide accurate information about CIPC, SARS, CIDB, CSD processes
- Guide users step-by-step through registration
- Be friendly, professional, and helpful
- Use South African terminology
- If unsure, say so and offer to connect with a human agent

Context:
{context}

Conversation History:
{history}

User: {question}
Assistant:`;

// Create prompt template
const prompt = PromptTemplate.fromTemplate(systemPrompt);

// Create chain
const chain = RunnableSequence.from([
  prompt,
  llm,
  new StringOutputParser(),
]);

// Generate response
export async function generateResponse(
  question: string,
  context: string,
  history: string
): Promise<string> {
  const response = await chain.invoke({
    question,
    context,
    history,
  });
  
  return response;
}
```

---

### RAG Implementation

```typescript
import { OpenAIEmbeddings } from "@langchain/openai";
import { PineconeStore } from "@langchain/community/vectorstores/pinecone";
import { Pinecone } from "@pinecone-database/pinecone";

// Initialize Pinecone
const pinecone = new Pinecone({
  apiKey: process.env.PINECONE_API_KEY!,
});

const index = pinecone.Index("comply360-knowledge");

// Create embeddings
const embeddings = new OpenAIEmbeddings({
  openAIApiKey: process.env.OPENAI_API_KEY,
});

// Initialize vector store
const vectorStore = await PineconeStore.fromExistingIndex(
  embeddings,
  { pineconeIndex: index }
);

// Retrieve relevant documents
export async function retrieveContext(query: string): Promise<string> {
  const results = await vectorStore.similaritySearch(query, 5);
  
  return results
    .map((doc) => doc.pageContent)
    .join("\n\n");
}

// Index new documents
export async function indexDocument(
  content: string,
  metadata: Record<string, any>
): Promise<void> {
  await vectorStore.addDocuments([
    {
      pageContent: content,
      metadata,
    },
  ]);
}
```

---

### WebSocket Handler

```typescript
import { Server as SocketIOServer } from "socket.io";
import { Server as HTTPServer } from "http";
import { verifyToken } from "./auth";
import { ConversationManager } from "./conversations/manager";
import { generateResponse } from "./llm/openai";
import { retrieveContext } from "./knowledge/retrieval";

export function initializeWebSocket(httpServer: HTTPServer) {
  const io = new SocketIOServer(httpServer, {
    cors: {
      origin: process.env.FRONTEND_URL,
      credentials: true,
    },
  });

  // Authentication middleware
  io.use(async (socket, next) => {
    const token = socket.handshake.auth.token;
    
    if (!token) {
      return next(new Error("Authentication required"));
    }
    
    try {
      const user = await verifyToken(token);
      socket.data.user = user;
      next();
    } catch (error) {
      next(new Error("Invalid token"));
    }
  });

  // Connection handler
  io.on("connection", (socket) => {
    console.log(`User connected: ${socket.data.user.id}`);

    const conversationManager = new ConversationManager(
      socket.data.user.id,
      socket.data.user.tenantId
    );

    // Message handler
    socket.on("message", async (data: { message: string }) => {
      try {
        // Add to conversation history
        await conversationManager.addMessage("user", data.message);

        // Get conversation context
        const history = await conversationManager.getHistory();
        const contextStr = conversationManager.formatHistory(history);

        // Retrieve relevant knowledge
        const knowledgeContext = await retrieveContext(data.message);

        // Generate response
        const response = await generateResponse(
          data.message,
          knowledgeContext,
          contextStr
        );

        // Add assistant response to history
        await conversationManager.addMessage("assistant", response);

        // Send response
        socket.emit("message", {
          message: response,
          timestamp: new Date(),
        });

        // Track analytics
        await trackMessage(socket.data.user.id, data.message, response);
      } catch (error) {
        console.error("Error processing message:", error);
        socket.emit("error", { message: "Failed to process message" });
      }
    });

    // Typing indicator
    socket.on("typing", () => {
      socket.broadcast.emit("typing", { userId: socket.data.user.id });
    });

    // Disconnect handler
    socket.on("disconnect", () => {
      console.log(`User disconnected: ${socket.data.user.id}`);
    });
  });

  return io;
}
```

---

### Conversation Manager

```typescript
import { Redis } from "ioredis";
import { db } from "../database";

interface Message {
  role: "user" | "assistant" | "system";
  content: string;
  timestamp: Date;
}

export class ConversationManager {
  private redis: Redis;
  private conversationId: string;

  constructor(
    private userId: string,
    private tenantId: string
  ) {
    this.redis = new Redis(process.env.REDIS_URL!);
    this.conversationId = `conversation:${userId}:${Date.now()}`;
  }

  async addMessage(role: "user" | "assistant", content: string): Promise<void> {
    const message: Message = {
      role,
      content,
      timestamp: new Date(),
    };

    // Store in Redis (for quick access)
    await this.redis.lpush(
      this.conversationId,
      JSON.stringify(message)
    );

    // Set expiry (24 hours)
    await this.redis.expire(this.conversationId, 86400);

    // Store in database (for analytics)
    await db.query(
      `INSERT INTO chat_messages 
       (conversation_id, user_id, tenant_id, role, content, created_at)
       VALUES ($1, $2, $3, $4, $5, $6)`,
      [
        this.conversationId,
        this.userId,
        this.tenantId,
        role,
        content,
        message.timestamp,
      ]
    );
  }

  async getHistory(limit: number = 10): Promise<Message[]> {
    const messages = await this.redis.lrange(this.conversationId, 0, limit - 1);
    
    return messages
      .map((msg) => JSON.parse(msg))
      .reverse();
  }

  formatHistory(messages: Message[]): string {
    return messages
      .map((msg) => `${msg.role}: ${msg.content}`)
      .join("\n");
  }

  async clearHistory(): Promise<void> {
    await this.redis.del(this.conversationId);
  }
}
```

---

### Intent Classification

```typescript
import { ChatOpenAI } from "@langchain/openai";
import { StructuredOutputParser } from "langchain/output_parsers";

interface Intent {
  name: string;
  confidence: number;
  entities: Record<string, string>;
}

const intentParser = StructuredOutputParser.fromNamesAndDescriptions({
  intent: "The detected intent (greeting, question, complaint, request_status, book_appointment)",
  confidence: "Confidence score from 0 to 1",
  entities: "Extracted entities as JSON object",
});

export async function classifyIntent(message: string): Promise<Intent> {
  const llm = new ChatOpenAI({ temperature: 0 });

  const prompt = `Classify the user's intent and extract entities.

User message: "${message}"

Possible intents:
- greeting: User is greeting or introducing themselves
- question: User is asking a question
- complaint: User is expressing dissatisfaction
- request_status: User wants to check status of application
- book_appointment: User wants to book an appointment
- help: User needs assistance
- goodbye: User is ending conversation

Extract entities like:
- company_name
- registration_number
- service_type
- date
- time

${intentParser.getFormatInstructions()}`;

  const response = await llm.predict(prompt);
  const parsed = await intentParser.parse(response);

  return {
    name: parsed.intent,
    confidence: parseFloat(parsed.confidence),
    entities: JSON.parse(parsed.entities || "{}"),
  };
}
```

---

### Actions Engine

```typescript
export interface Action {
  execute(params: any): Promise<ActionResult>;
}

export interface ActionResult {
  success: boolean;
  data?: any;
  message: string;
}

// Check registration status
export class CheckStatusAction implements Action {
  async execute(params: { registrationId: string }): Promise<ActionResult> {
    const registration = await db.query(
      "SELECT * FROM registrations WHERE id = $1",
      [params.registrationId]
    );

    if (!registration.rows[0]) {
      return {
        success: false,
        message: "Registration not found",
      };
    }

    const reg = registration.rows[0];

    return {
      success: true,
      data: reg,
      message: `Your registration (${reg.company_name}) is currently ${reg.status}. 
                Estimated completion: ${reg.estimated_completion}`,
    };
  }
}

// Calculate fee
export class CalculateFeeAction implements Action {
  async execute(params: { serviceType: string }): Promise<ActionResult> {
    const fees = {
      company_registration: 175,
      tax_clearance: 0,
      vat_registration: 0,
      cidb_registration: 1000,
      // ... more services
    };

    const fee = fees[params.serviceType];

    if (fee === undefined) {
      return {
        success: false,
        message: "Service type not found",
      };
    }

    return {
      success: true,
      data: { fee },
      message: `The fee for ${params.serviceType} is R${fee}`,
    };
  }
}

// Transfer to human agent
export class TransferAgentAction implements Action {
  async execute(params: { conversationId: string }): Promise<ActionResult> {
    // Find available agent
    const agent = await findAvailableAgent();

    if (!agent) {
      return {
        success: false,
        message: "No agents currently available. We'll email you shortly.",
      };
    }

    // Assign conversation to agent
    await assignConversation(params.conversationId, agent.id);

    return {
      success: true,
      data: { agent },
      message: `Connecting you with ${agent.name}...`,
    };
  }
}
```

---

## Chat Widget (Frontend)

```typescript
// ChatWidget.svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import io, { Socket } from 'socket.io-client';
  
  let messages: Array<{role: string, content: string, timestamp: Date}> = [];
  let inputMessage = '';
  let socket: Socket;
  let isOpen = false;
  let isTyping = false;
  
  onMount(() => {
    socket = io(import.meta.env.VITE_CHATBOT_URL, {
      auth: {
        token: getAuthToken(),
      },
    });
    
    socket.on('message', (data) => {
      messages = [...messages, {
        role: 'assistant',
        content: data.message,
        timestamp: new Date(data.timestamp),
      }];
      isTyping = false;
    });
    
    socket.on('typing', () => {
      isTyping = true;
    });
  });
  
  function sendMessage() {
    if (!inputMessage.trim()) return;
    
    messages = [...messages, {
      role: 'user',
      content: inputMessage,
      timestamp: new Date(),
    }];
    
    socket.emit('message', { message: inputMessage });
    inputMessage = '';
    isTyping = true;
  }
  
  function toggleChat() {
    isOpen = !isOpen;
  }
</script>

<div class="chatbot-container" class:open={isOpen}>
  {#if isOpen}
    <div class="chat-header">
      <h3>Comply360 Assistant</h3>
      <button on:click={toggleChat}>×</button>
    </div>
    
    <div class="chat-messages">
      {#each messages as message}
        <div class="message {message.role}">
          <div class="content">{message.content}</div>
          <div class="timestamp">
            {message.timestamp.toLocaleTimeString()}
          </div>
        </div>
      {/each}
      
      {#if isTyping}
        <div class="message assistant">
          <div class="typing-indicator">
            <span></span><span></span><span></span>
          </div>
        </div>
      {/if}
    </div>
    
    <div class="chat-input">
      <input
        type="text"
        bind:value={inputMessage}
        on:keypress={(e) => e.key === 'Enter' && sendMessage()}
        placeholder="Type your message..."
      />
      <button on:click={sendMessage}>Send</button>
    </div>
  {:else}
    <button class="chat-toggle" on:click={toggleChat}>
      💬 Chat with us
    </button>
  {/if}
</div>

<style>
  .chatbot-container {
    position: fixed;
    bottom: 20px;
    right: 20px;
    z-index: 1000;
  }
  
  .chatbot-container.open {
    width: 400px;
    height: 600px;
    background: white;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0,0,0,0.15);
    display: flex;
    flex-direction: column;
  }
  
  .chat-header {
    padding: 16px;
    background: #0066CC;
    color: white;
    border-radius: 12px 12px 0 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .chat-messages {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }
  
  .message {
    margin-bottom: 12px;
    padding: 12px;
    border-radius: 8px;
    max-width: 80%;
  }
  
  .message.user {
    background: #0066CC;
    color: white;
    margin-left: auto;
  }
  
  .message.assistant {
    background: #F5F5F5;
    color: #333;
  }
  
  .chat-input {
    padding: 16px;
    border-top: 1px solid #E0E0E0;
    display: flex;
    gap: 8px;
  }
  
  .chat-input input {
    flex: 1;
    padding: 12px;
    border: 1px solid #E0E0E0;
    border-radius: 8px;
  }
  
  .chat-toggle {
    padding: 16px 24px;
    background: #0066CC;
    color: white;
    border: none;
    border-radius: 50px;
    cursor: pointer;
    font-size: 16px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  }
  
  .typing-indicator span {
    display: inline-block;
    width: 8px;
    height: 8px;
    background: #999;
    border-radius: 50%;
    margin: 0 2px;
    animation: typing 1.4s infinite;
  }
  
  @keyframes typing {
    0%, 60%, 100% { transform: translateY(0); }
    30% { transform: translateY(-10px); }
  }
</style>
```

---

## Database Schema

```sql
CREATE TABLE chat_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    assigned_agent_id UUID REFERENCES users(id),
    rating INTEGER,
    feedback TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES chat_conversations(id),
    user_id UUID NOT NULL REFERENCES users(id),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    role VARCHAR(20) NOT NULL, -- 'user' | 'assistant' | 'system'
    content TEXT NOT NULL,
    intent VARCHAR(50),
    confidence DECIMAL(5,4),
    entities JSONB,
    action_taken VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE chat_knowledge_base (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    keywords TEXT[],
    embedding VECTOR(1536), -- for similarity search
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chat_conversations_user ON chat_conversations(user_id);
CREATE INDEX idx_chat_messages_conversation ON chat_messages(conversation_id);
CREATE INDEX idx_chat_knowledge_category ON chat_knowledge_base(category);
```

---

## Performance Requirements

- **Response Time**: < 2 seconds for simple queries
- **Response Time (RAG)**: < 4 seconds for knowledge-based queries
- **Concurrent Users**: Support 1000+ simultaneous conversations
- **Uptime**: 99.9%
- **Language Detection**: < 100ms
- **WebSocket Latency**: < 100ms

---

## Cost Estimation

**Per 1000 Messages:**
- OpenAI API (GPT-4 Turbo): $0.30 (input) + $0.60 (output) = $0.90
- Pinecone (Vector DB): $70/month (100k vectors)
- Redis: $20/month
- **Total: ~$1/1000 messages**

**Monthly Costs (10k users, 5 messages/user):**
- OpenAI: $45
- Pinecone: $70
- Redis: $20
- Compute: $100
- **Total: ~$235/month**

---

## Success Metrics

- Response accuracy: > 90%
- User satisfaction (CSAT): > 4.5/5
- Support ticket reduction: > 60%
- First response resolution: > 75%
- Average response time: < 3 seconds
- Handoff rate: < 15%

---

**Next Steps:** See `tasks.md` for implementation tasks

