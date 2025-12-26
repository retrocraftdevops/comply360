#!/usr/bin/env node

/**
 * Comply360 AI Code Validation System
 * 
 * Validates code against enterprise standards:
 * - Security: Tenant isolation, input validation, error handling
 * - TypeScript: Proper types, no any types
 * - Performance: Memoization, optimization
 * - Standards: No console logs, documentation
 * - Database: Proper indexes, where clauses
 * - API: Response format, HTTP status codes
 * - UI: Accessibility, loading states
 * 
 * Usage:
 *   npm run ai:validate path/to/file.ts
 *   npm run ai:validate:all
 *   npm run ai:validate:all --categories "security,typescript"
 */

const fs = require('fs');
const path = require('path');

// Validation rules by category
const VALIDATION_RULES = {
  security: [
    {
      name: 'Tenant Isolation',
      pattern: /where.*tenantId/gi,
      severity: 'critical',
      message: 'Always filter database queries by tenantId for multi-tenant isolation'
    },
    {
      name: 'Input Validation',
      pattern: /(\.parse\(|\.safeParse\(|validator\.)/gi,
      severity: 'high',
      message: 'Use Zod or validator for input validation'
    },
    {
      name: 'Error Handling',
      pattern: /(try\s*{|catch\s*\()/gi,
      severity: 'high',
      message: 'Proper try-catch error handling required'
    },
    {
      name: 'No Hardcoded Secrets',
      pattern: /(password|secret|api_key)\s*=\s*['"]/gi,
      severity: 'critical',
      message: 'Never hardcode secrets, use environment variables'
    }
  ],
  
  typescript: [
    {
      name: 'No Any Types',
      pattern: /:\s*any(\s|;|,|\))/gi,
      severity: 'high',
      message: 'Avoid using "any" type, use specific types'
    },
    {
      name: 'Type Definitions',
      pattern: /(interface|type)\s+\w+/gi,
      severity: 'medium',
      message: 'Good: Type definitions present'
    },
    {
      name: 'Strict Null Checks',
      pattern: /\?\./gi,
      severity: 'low',
      message: 'Good: Using optional chaining'
    }
  ],
  
  performance: [
    {
      name: 'Memoization',
      pattern: /(useMemo|useCallback|React\.memo)/gi,
      severity: 'medium',
      message: 'Good: Using memoization for performance'
    },
    {
      name: 'Lazy Loading',
      pattern: /(React\.lazy|dynamic\()/gi,
      severity: 'low',
      message: 'Good: Using lazy loading'
    }
  ],
  
  standards: [
    {
      name: 'No Console Logs',
      pattern: /console\.(log|warn|error)/gi,
      severity: 'medium',
      message: 'Remove console logs before production'
    },
    {
      name: 'Documentation',
      pattern: /\/\*\*[\s\S]*?\*\//g,
      severity: 'low',
      message: 'Good: JSDoc comments present'
    }
  ],
  
  database: [
    {
      name: 'Proper Indexes',
      pattern: /@@index/gi,
      severity: 'medium',
      message: 'Good: Database indexes defined'
    },
    {
      name: 'Where Clauses',
      pattern: /where:\s*{/gi,
      severity: 'high',
      message: 'Good: Using where clauses for filtering'
    }
  ],
  
  api: [
    {
      name: 'Response Format',
      pattern: /Response\.json\(/gi,
      severity: 'high',
      message: 'Good: Using proper response format'
    },
    {
      name: 'HTTP Status Codes',
      pattern: /status:\s*\d{3}/gi,
      severity: 'high',
      message: 'Good: Proper HTTP status codes'
    }
  ],
  
  ui: [
    {
      name: 'Accessibility',
      pattern: /(aria-label|aria-|role=)/gi,
      severity: 'high',
      message: 'Good: Accessibility attributes present'
    },
    {
      name: 'Loading States',
      pattern: /(isLoading|loading|isPending)/gi,
      severity: 'medium',
      message: 'Good: Loading states handled'
    }
  ]
};

// Scan file and apply validation rules
function validateFile(filePath, categories = Object.keys(VALIDATION_RULES)) {
  if (!fs.existsSync(filePath)) {
    console.error(`❌ File not found: ${filePath}`);
    return null;
  }

  const content = fs.readFileSync(filePath, 'utf8');
  const results = {
    filePath,
    passed: [],
    failed: [],
    score: 0
  };

  // Apply rules from selected categories
  for (const category of categories) {
    const rules = VALIDATION_RULES[category] || [];
    
    for (const rule of rules) {
      const matches = content.match(rule.pattern);
      const found = matches && matches.length > 0;

      const result = {
        category,
        name: rule.name,
        severity: rule.severity,
        message: rule.message,
        found
      };

      // Positive rules (should have)
      if (rule.message.startsWith('Good:')) {
        if (found) {
          results.passed.push(result);
        }
      } 
      // Negative rules (should NOT have)
      else {
        if (!found) {
          results.passed.push({...result, message: `✓ ${rule.message} check passed`});
        } else {
          results.failed.push({...result, count: matches.length});
        }
      }
    }
  }

  // Calculate score
  const total = results.passed.length + results.failed.length;
  results.score = total > 0 ? Math.round((results.passed.length / total) * 100) : 0;

  return results;
}

// Print validation results
function printResults(results, verbose = false) {
  console.log(`\n📄 File: ${results.filePath}`);
  console.log(`📊 Score: ${results.score}% (${results.passed.length}/${results.passed.length + results.failed.length} checks passed)`);

  // Critical failures
  const critical = results.failed.filter(r => r.severity === 'critical');
  if (critical.length > 0) {
    console.log(`\n🚨 CRITICAL ISSUES (${critical.length}):`);
    critical.forEach(r => {
      console.log(`   ❌ [${r.category}] ${r.name}: ${r.message}`);
    });
  }

  // High severity failures
  const high = results.failed.filter(r => r.severity === 'high');
  if (high.length > 0) {
    console.log(`\n⚠️  HIGH PRIORITY (${high.length}):`);
    high.length(r => {
      console.log(`   ❌ [${r.category}] ${r.name}: ${r.message}`);
    });
  }

  // Medium/Low severity (only in verbose mode)
  if (verbose) {
    const medium = results.failed.filter(r => r.severity === 'medium');
    const low = results.failed.filter(r => r.severity === 'low');
    
    if (medium.length > 0) {
      console.log(`\n📋 MEDIUM PRIORITY (${medium.length}):`);
      medium.forEach(r => {
        console.log(`   ⚠️  [${r.category}] ${r.name}: ${r.message}`);
      });
    }
    
    if (low.length > 0) {
      console.log(`\n💡 SUGGESTIONS (${low.length}):`);
      low.forEach(r => {
        console.log(`   ℹ️  [${r.category}] ${r.name}: ${r.message}`);
      });
    }
  }

  // Pass/Fail determination
  const hasCritical = critical.length > 0;
  const scorePass = results.score >= 80;
  
  if (hasCritical) {
    console.log(`\n❌ VALIDATION FAILED: Critical issues must be fixed\n`);
    return false;
  } else if (!scorePass) {
    console.log(`\n⚠️  VALIDATION WARNING: Score below 80% (current: ${results.score}%)\n`);
    return false;
  } else {
    console.log(`\n✅ VALIDATION PASSED: Score ${results.score}% (>= 80%)\n`);
    return true;
  }
}

// Main execution
function main() {
  const args = process.argv.slice(2);
  
  // Parse arguments
  let filePath = null;
  let validateAll = false;
  let verbose = false;
  let categories = Object.keys(VALIDATION_RULES);

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    
    if (arg === '--all') {
      validateAll = true;
    } else if (arg === '--verbose') {
      verbose = true;
    } else if (arg === '--categories') {
      categories = args[i + 1].split(',');
      i++;
    } else if (arg === '--help') {
      console.log(`
Comply360 AI Code Validation System

Usage:
  npm run ai:validate <file>                         # Validate single file
  npm run ai:validate:all                            # Validate all files
  npm run ai:validate:all --categories "security"    # Validate specific categories
  npm run ai:validate:verbose <file>                 # Verbose output

Categories:
  security, typescript, performance, standards, database, api, ui

Examples:
  npm run ai:validate apps/web/app/page.tsx
  npm run ai:validate:all --categories "security,typescript"
  npm run ai:validate:verbose apps/api/internal/handlers/registration.go
      `);
      process.exit(0);
    } else if (!arg.startsWith('--')) {
      filePath = arg;
    }
  }

  // Validate files
  if (validateAll) {
    console.log('🔍 Validating all files...\n');
    
    // Find all TypeScript/JavaScript files
    const glob = require('glob');
    const files = glob.sync('**/*.{ts,tsx,js,jsx}', {
      ignore: ['node_modules/**', 'dist/**', 'build/**', '.next/**']
    });

    let totalScore = 0;
    let failedFiles = [];

    files.forEach(file => {
      const results = validateFile(file, categories);
      if (results) {
        const passed = printResults(results, verbose);
        totalScore += results.score;
        if (!passed) {
          failedFiles.push({ file, score: results.score });
        }
      }
    });

    const avgScore = Math.round(totalScore / files.length);
    
    console.log(`\n${'='.repeat(60)}`);
    console.log(`📊 SUMMARY: ${files.length} files validated`);
    console.log(`📈 Average Score: ${avgScore}%`);
    console.log(`✅ Passed: ${files.length - failedFiles.length}`);
    console.log(`❌ Failed: ${failedFiles.length}`);
    
    if (failedFiles.length > 0) {
      console.log(`\n⚠️  Files needing attention:`);
      failedFiles.forEach(({ file, score }) => {
        console.log(`   ${file} (${score}%)`);
      });
    }
    
    console.log(`${'='.repeat(60)}\n`);
    
    process.exit(failedFiles.length > 0 ? 1 : 0);
    
  } else if (filePath) {
    const results = validateFile(filePath, categories);
    if (results) {
      const passed = printResults(results, verbose);
      process.exit(passed ? 0 : 1);
    } else {
      process.exit(1);
    }
  } else {
    console.error('❌ Please specify a file path or use --all flag');
    console.log('Run --help for usage information');
    process.exit(1);
  }
}

// Install glob if not available
try {
  require('glob');
  main();
} catch (e) {
  console.log('📦 Installing required dependencies...');
  require('child_process').execSync('npm install glob', { stdio: 'inherit' });
  main();
}
