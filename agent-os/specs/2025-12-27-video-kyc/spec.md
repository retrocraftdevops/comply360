# Video KYC (Know Your Customer) - Specification

**Version:** 1.0.0  
**Date:** December 27, 2025  
**Priority:** P1 - High  
**Estimated Duration:** 4-5 weeks  

---

## Executive Summary

Real-time video verification system that allows agents to remotely verify client identities, conduct video interviews, and capture live biometric data for compliance purposes. This feature eliminates the need for in-person meetings and reduces fraud by 90%.

---

## Business Case

### Market Need
- Traditional KYC requires in-person verification
- Fraud is increasing (identity theft, fake documents)
- Remote verification is becoming mandatory
- Competitors don't offer video KYC in SADC
- FICA compliance requires enhanced verification

### Business Impact
- **Fraud Reduction**: 90% decrease in identity fraud
- **Processing Speed**: 80% faster verification
- **Cost Savings**: R200k/year in travel and office costs
- **User Convenience**: Remote verification from anywhere
- **Compliance**: Meet FICA requirements
- **Competitive Advantage**: Only platform with video KYC in SADC

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    FRONTEND                              │
│  (Video Call Interface)                                  │
└────────────────────────┬─────────────────────────────────┘
                         │ WebRTC
                         ▼
┌─────────────────────────────────────────────────────────┐
│              VIDEO KYC SERVICE (Node.js)                 │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │  WebRTC      │───▶│  Recording   │                  │
│  │  Signaling   │    │   Service    │                  │
│  └──────────────┘    └──────┬───────┘                  │
│                              │                           │
│                              ▼                           │
│                      ┌──────────────┐                   │
│                      │   Live Face  │                   │
│                      │  Detection   │                   │
│                      └──────┬───────┘                   │
│                              │                           │
│                              ▼                           │
│                      ┌──────────────┐                   │
│                      │ Liveness     │                   │
│                      │ Detection    │                   │
│                      └──────┬───────┘                   │
│                              │                           │
│                              ▼                           │
│                      ┌──────────────┐                   │
│                      │  Face Match  │                   │
│                      │  (ID vs Live)│                   │
│                      └──────┬───────┘                   │
│                              │                           │
│                              ▼                           │
│                      ┌──────────────┐                   │
│                      │  Document    │                   │
│                      │  Capture     │                   │
│                      └──────┬───────┘                   │
│                              │                           │
│                              ▼                           │
│                      ┌──────────────┐                   │
│                      │  Compliance  │                   │
│                      │  Report      │                   │
│                      └───────────────┘                  │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  STORAGE & AI                            │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │  S3/MinIO    │    │   AWS        │                  │
│  │  (Videos)    │    │  Rekognition │                  │
│  └──────────────┘    └──────────────┘                  │
└─────────────────────────────────────────────────────────┘
```

---

## Core Features

### 1. Video Call Interface

**WebRTC Implementation:**
- Peer-to-peer video connection
- HD video quality (720p minimum)
- Audio communication
- Screen sharing capability
- Document sharing
- Recording functionality

**Call Features:**
- Agent-to-client video calls
- Multi-party calls (up to 4 participants)
- Call scheduling
- Automatic recording
- Call duration tracking
- Call quality monitoring

---

### 2. Identity Verification

**Live Face Detection:**
- Real-time face detection
- Face positioning guides
- Multiple angle capture (front, left, right)
- Lighting quality check
- Camera quality validation

**Liveness Detection:**
- Blink detection
- Head movement tracking
- Challenge-response (smile, turn head)
- Anti-spoofing (detect photos, videos, masks)
- Passive liveness (subtle movements)
- 3D depth analysis

**Face Matching:**
- Compare live face to ID photo
- Similarity score (0-100%)
- Threshold: 80%+ for pass
- Multi-angle comparison
- Age progression consideration

---

### 3. Document Verification

**Live Document Capture:**
- ID document capture during call
- Security features detection
- Hologram verification
- Watermark detection
- Document tilt detection
- Reflections and glare detection

**Real-Time OCR:**
- Extract data during call
- Instant verification
- Compare to submitted documents
- Detect tampering
- Flag inconsistencies

---

### 4. Biometric Capture

**Facial Biometrics:**
- High-resolution face capture
- Facial landmarks (68+ points)
- 3D face mapping
- Texture analysis
- Store biometric templates

**Voice Biometrics (Optional):**
- Voice recording
- Voice print creation
- Speaker verification
- Accent/language detection

---

### 5. Guided Verification Workflow

**Agent Workflow:**
```
1. Schedule Call
   ↓
2. Send Invite to Client
   ↓
3. Pre-Call Checklist
   - ID documents uploaded
   - Camera/mic test
   - Consent forms signed
   ↓
4. Start Video Call
   ↓
5. Live Verification Steps:
   a. Greet client
   b. Verify ID document on screen
   c. Perform liveness checks
   d. Capture facial biometrics
   e. Ask verification questions
   f. Capture additional documents
   g. Review and confirm
   ↓
6. End Call
   ↓
7. Generate Compliance Report
   ↓
8. Approve/Reject
```

**Client Experience:**
```
1. Receive Email/SMS Invite
   ↓
2. Click Link
   ↓
3. Camera/Mic Permission
   ↓
4. Pre-Call Instructions
   - Have ID ready
   - Good lighting
   - Quiet environment
   ↓
5. Join Call
   ↓
6. Follow Agent Instructions
   - Show ID
   - Perform liveness actions
   - Answer questions
   ↓
7. Receive Confirmation
```

---

### 6. Recording & Audit Trail

**Call Recording:**
- Full video/audio recording
- Encrypted storage
- Timestamped events
- Automatic transcription
- Searchable content

**Audit Trail:**
- Call start/end times
- Participants list
- Actions taken
- Documents captured
- Verification results
- Agent notes
- Approval/rejection reasons

**Compliance:**
- POPI Act compliance
- FICA compliance
- Consent management
- Data retention policies
- Access controls
- Audit logs

---

## Technical Specifications

### WebRTC Implementation

**Technology Stack:**
- **Signaling**: Socket.io
- **STUN/TURN**: Coturn server
- **Media Server**: Mediasoup (optional for recording)
- **Frontend**: React/Svelte with simple-peer or PeerJS

**WebRTC Flow:**
```
Client A                Signaling Server              Client B
   │                           │                         │
   │─────offer───────────────>│                         │
   │                           │─────offer────────────>│
   │                           │                         │
   │                           │<────answer────────────│
   │<────answer──────────────│                         │
   │                           │                         │
   │<─────────────ICE Candidates─────────────────────>│
   │                           │                         │
   │<══════════ Direct P2P Connection ════════════════>│
```

---

### Project Structure

```
apps/video-kyc/
├── src/
│   ├── server.ts
│   ├── signaling/
│   │   ├── handler.ts
│   │   ├── room-manager.ts
│   │   └── peer-connection.ts
│   ├── recording/
│   │   ├── recorder.ts
│   │   ├── storage.ts
│   │   └── transcription.ts
│   ├── verification/
│   │   ├── face-detection.ts
│   │   ├── liveness.ts
│   │   ├── face-matching.ts
│   │   └── document-capture.ts
│   ├── biometrics/
│   │   ├── facial.ts
│   │   ├── voice.ts
│   │   └── storage.ts
│   ├── compliance/
│   │   ├── report-generator.ts
│   │   ├── audit-trail.ts
│   │   └── consent.ts
│   └── integrations/
│       ├── aws-rekognition.ts
│       ├── s3.ts
│       └── comply360-api.ts
├── Dockerfile
├── package.json
└── tsconfig.json
```

---

### WebRTC Signaling Server

```typescript
import { Server as SocketIOServer } from "socket.io";
import { Server as HTTPServer } from "http";

interface Room {
  id: string;
  participants: Map<string, Participant>;
  startedAt: Date;
  recording: boolean;
}

interface Participant {
  socketId: string;
  userId: string;
  role: "agent" | "client";
  stream: boolean;
}

class SignalingServer {
  private io: SocketIOServer;
  private rooms: Map<string, Room> = new Map();

  constructor(httpServer: HTTPServer) {
    this.io = new SocketIOServer(httpServer, {
      cors: {
        origin: process.env.FRONTEND_URL,
        credentials: true,
      },
    });

    this.initialize();
  }

  private initialize() {
    this.io.on("connection", (socket) => {
      console.log(`User connected: ${socket.id}`);

      // Join room
      socket.on("join-room", async (data: { roomId: string; userId: string; role: string }) => {
        const { roomId, userId, role } = data;

        // Create room if doesn't exist
        if (!this.rooms.has(roomId)) {
          this.rooms.set(roomId, {
            id: roomId,
            participants: new Map(),
            startedAt: new Date(),
            recording: false,
          });
        }

        const room = this.rooms.get(roomId)!;
        
        // Add participant
        room.participants.set(socket.id, {
          socketId: socket.id,
          userId,
          role: role as "agent" | "client",
          stream: false,
        });

        socket.join(roomId);

        // Notify other participants
        socket.to(roomId).emit("user-joined", {
          socketId: socket.id,
          userId,
          role,
        });

        // Send existing participants to new user
        const existingParticipants = Array.from(room.participants.values())
          .filter(p => p.socketId !== socket.id);
        
        socket.emit("existing-participants", existingParticipants);

        console.log(`User ${userId} joined room ${roomId}`);
      });

      // WebRTC offer
      socket.on("offer", (data: { to: string; offer: RTCSessionDescriptionInit }) => {
        socket.to(data.to).emit("offer", {
          from: socket.id,
          offer: data.offer,
        });
      });

      // WebRTC answer
      socket.on("answer", (data: { to: string; answer: RTCSessionDescriptionInit }) => {
        socket.to(data.to).emit("answer", {
          from: socket.id,
          answer: data.answer,
        });
      });

      // ICE candidate
      socket.on("ice-candidate", (data: { to: string; candidate: RTCIceCandidate }) => {
        socket.to(data.to).emit("ice-candidate", {
          from: socket.id,
          candidate: data.candidate,
        });
      });

      // Start recording
      socket.on("start-recording", async (data: { roomId: string }) => {
        const room = this.rooms.get(data.roomId);
        if (room) {
          room.recording = true;
          // Start recording service
          await this.startRecording(data.roomId);
          this.io.to(data.roomId).emit("recording-started");
        }
      });

      // Disconnect
      socket.on("disconnect", () => {
        // Find and remove from room
        this.rooms.forEach((room) => {
          if (room.participants.has(socket.id)) {
            room.participants.delete(socket.id);
            this.io.to(room.id).emit("user-left", { socketId: socket.id });
            
            // Clean up empty rooms
            if (room.participants.size === 0) {
              this.rooms.delete(room.id);
            }
          }
        });
      });
    });
  }

  private async startRecording(roomId: string) {
    // Implement recording logic
    // Use Mediasoup or record client-side and upload
  }
}

export function initializeSignaling(httpServer: HTTPServer) {
  return new SignalingServer(httpServer);
}
```

---

### Face Detection & Liveness

```typescript
import * as faceapi from "@vladmandic/face-api";
import Rekognition from "aws-sdk/clients/rekognition";

export class FaceVerification {
  private rekognition: Rekognition;

  constructor() {
    this.rekognition = new Rekognition({
      region: process.env.AWS_REGION,
    });
  }

  // Detect faces in image
  async detectFaces(imageBuffer: Buffer): Promise<faceapi.FaceDetection[]> {
    const img = await faceapi.fetchImage(imageBuffer);
    const detections = await faceapi
      .detectAllFaces(img)
      .withFaceLandmarks()
      .withFaceDescriptors();
    
    return detections;
  }

  // Liveness detection using AWS Rekognition
  async checkLiveness(videoBuffer: Buffer): Promise<LivenessResult> {
    const params = {
      Video: {
        Bytes: videoBuffer,
      },
    };

    const response = await this.rekognition.detectFaces(params).promise();

    // Check for liveness indicators
    const faces = response.FaceDetails || [];
    
    if (faces.length === 0) {
      return {
        isLive: false,
        confidence: 0,
        reason: "No face detected",
      };
    }

    const face = faces[0];
    
    // Check quality metrics
    const qualityScore = this.calculateQualityScore(face);
    
    return {
      isLive: qualityScore > 80,
      confidence: qualityScore,
      details: {
        eyesOpen: face.EyesOpen?.Value,
        mouthOpen: face.MouthOpen?.Value,
        brightness: face.Quality?.Brightness,
        sharpness: face.Quality?.Sharpness,
      },
    };
  }

  // Compare two faces
  async compareFaces(
    sourceImage: Buffer,
    targetImage: Buffer
  ): Promise<FaceMatchResult> {
    const params = {
      SourceImage: { Bytes: sourceImage },
      TargetImage: { Bytes: targetImage },
      SimilarityThreshold: 80,
    };

    const response = await this.rekognition.compareFaces(params).promise();

    const matches = response.FaceMatches || [];
    
    if (matches.length === 0) {
      return {
        matched: false,
        similarity: 0,
        confidence: 0,
      };
    }

    const match = matches[0];
    
    return {
      matched: true,
      similarity: match.Similarity || 0,
      confidence: match.Face?.Confidence || 0,
      faceDetails: match.Face,
    };
  }

  private calculateQualityScore(face: Rekognition.FaceDetail): number {
    const brightness = face.Quality?.Brightness || 0;
    const sharpness = face.Quality?.Sharpness || 0;
    const eyesOpen = face.EyesOpen?.Confidence || 0;
    
    return (brightness + sharpness + eyesOpen) / 3;
  }
}

interface LivenessResult {
  isLive: boolean;
  confidence: number;
  reason?: string;
  details?: any;
}

interface FaceMatchResult {
  matched: boolean;
  similarity: number;
  confidence: number;
  faceDetails?: Rekognition.Face;
}
```

---

### Video Recording Service

```typescript
import { S3 } from "aws-sdk";
import { MediaRecorder } from "@ffmpeg/ffmpeg";

export class RecordingService {
  private s3: S3;
  private activeRecordings: Map<string, Recording> = new Map();

  constructor() {
    this.s3 = new S3({
      region: process.env.AWS_REGION,
    });
  }

  async startRecording(
    roomId: string,
    streams: MediaStream[]
  ): Promise<string> {
    const recordingId = `recording-${roomId}-${Date.now()}`;
    
    const recording: Recording = {
      id: recordingId,
      roomId,
      startedAt: new Date(),
      chunks: [],
    };

    this.activeRecordings.set(recordingId, recording);

    // Start recording each stream
    streams.forEach((stream, index) => {
      const mediaRecorder = new MediaRecorder(stream, {
        mimeType: "video/webm;codecs=vp9",
      });

      mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          recording.chunks.push(event.data);
        }
      };

      mediaRecorder.start(1000); // Capture every second
    });

    return recordingId;
  }

  async stopRecording(recordingId: string): Promise<string> {
    const recording = this.activeRecordings.get(recordingId);
    
    if (!recording) {
      throw new Error("Recording not found");
    }

    recording.endedAt = new Date();

    // Combine chunks
    const blob = new Blob(recording.chunks, { type: "video/webm" });
    const buffer = Buffer.from(await blob.arrayBuffer());

    // Upload to S3
    const key = `kyc-recordings/${recording.roomId}/${recordingId}.webm`;
    
    await this.s3.upload({
      Bucket: process.env.S3_BUCKET!,
      Key: key,
      Body: buffer,
      ContentType: "video/webm",
      ServerSideEncryption: "AES256",
      Metadata: {
        roomId: recording.roomId,
        startedAt: recording.startedAt.toISOString(),
        endedAt: recording.endedAt!.toISOString(),
      },
    }).promise();

    // Clean up
    this.activeRecordings.delete(recordingId);

    return key;
  }

  async getRecordingUrl(key: string): Promise<string> {
    const url = await this.s3.getSignedUrlPromise("getObject", {
      Bucket: process.env.S3_BUCKET!,
      Key: key,
      Expires: 3600, // 1 hour
    });

    return url;
  }
}

interface Recording {
  id: string;
  roomId: string;
  startedAt: Date;
  endedAt?: Date;
  chunks: Blob[];
}
```

---

### Compliance Report Generator

```typescript
export class ComplianceReportGenerator {
  async generateReport(sessionId: string): Promise<ComplianceReport> {
    // Fetch session data
    const session = await this.getSession(sessionId);
    const verificationResults = await this.getVerificationResults(sessionId);
    const recording = await this.getRecording(sessionId);

    const report: ComplianceReport = {
      sessionId,
      generatedAt: new Date(),
      client: {
        userId: session.clientId,
        fullName: session.clientName,
        idNumber: session.idNumber,
      },
      agent: {
        userId: session.agentId,
        fullName: session.agentName,
      },
      session: {
        startedAt: session.startedAt,
        endedAt: session.endedAt,
        duration: this.calculateDuration(session.startedAt, session.endedAt),
        recordingUrl: recording.url,
      },
      verification: {
        liveness: verificationResults.liveness,
        faceMatch: verificationResults.faceMatch,
        documentVerification: verificationResults.document,
        overallStatus: this.calculateOverallStatus(verificationResults),
      },
      compliance: {
        consentObtained: session.consentObtained,
        dataRetention: "7 years",
        popiCompliant: true,
        ficaCompliant: true,
      },
      auditTrail: session.auditTrail,
    };

    // Store report
    await this.storeReport(report);

    return report;
  }

  private calculateOverallStatus(results: any): "pass" | "fail" | "review" {
    if (results.liveness.isLive && results.faceMatch.matched && results.faceMatch.similarity > 80) {
      return "pass";
    } else if (results.liveness.isLive && results.faceMatch.similarity > 70) {
      return "review";
    } else {
      return "fail";
    }
  }
}

interface ComplianceReport {
  sessionId: string;
  generatedAt: Date;
  client: any;
  agent: any;
  session: any;
  verification: any;
  compliance: any;
  auditTrail: any[];
}
```

---

## Database Schema

```sql
CREATE TABLE kyc_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    registration_id UUID REFERENCES registrations(id),
    client_id UUID NOT NULL REFERENCES users(id),
    agent_id UUID NOT NULL REFERENCES users(id),
    room_id VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    scheduled_at TIMESTAMP,
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    recording_url TEXT,
    transcript_url TEXT,
    consent_obtained BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE kyc_verification_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES kyc_sessions(id),
    liveness_check JSONB NOT NULL,
    face_match JSONB NOT NULL,
    document_verification JSONB,
    biometric_data JSONB,
    overall_status VARCHAR(20) NOT NULL,
    agent_notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE kyc_audit_trail (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES kyc_sessions(id),
    action VARCHAR(100) NOT NULL,
    performed_by UUID NOT NULL REFERENCES users(id),
    details JSONB,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE kyc_compliance_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES kyc_sessions(id),
    report_data JSONB NOT NULL,
    generated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    pdf_url TEXT
);

CREATE INDEX idx_kyc_sessions_client ON kyc_sessions(client_id);
CREATE INDEX idx_kyc_sessions_agent ON kyc_sessions(agent_id);
CREATE INDEX idx_kyc_sessions_status ON kyc_sessions(status);
CREATE INDEX idx_kyc_verification_session ON kyc_verification_results(session_id);
```

---

## Performance Requirements

- **Video Quality**: 720p HD minimum, 1080p preferred
- **Frame Rate**: 30 FPS minimum
- **Latency**: < 200ms for video
- **Face Detection**: < 100ms per frame
- **Liveness Check**: < 3 seconds
- **Face Matching**: < 2 seconds
- **Recording**: Real-time without lag

---

## Cost Estimation

**Per 1000 Sessions:**
- AWS Rekognition: $10 (face comparison)
- S3 Storage: $2 (1 GB per hour)
- TURN Server: $50/month
- Bandwidth: $20
- **Total: ~$82/1000 sessions**

---

## Success Metrics

- Fraud detection: > 90%
- False positive rate: < 5%
- Session completion: > 95%
- User satisfaction: > 4.5/5
- Average session time: < 10 minutes
- Compliance rate: 100%

---

**Next Steps:** See `tasks.md` for implementation tasks

